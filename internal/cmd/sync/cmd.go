// Package sync provides functionality around syncing Git branches.
package sync

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/flags"
	"github.com/git-town/git-town/v17/internal/cli/print"
	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/config"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/execute"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/hosting"
	"github.com/git-town/git-town/v17/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/undo/undoconfig"
	"github.com/git-town/git-town/v17/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v17/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v17/internal/vm/opcodes"
	"github.com/git-town/git-town/v17/internal/vm/optimizer"
	"github.com/git-town/git-town/v17/internal/vm/program"
	"github.com/git-town/git-town/v17/internal/vm/runstate"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/git-town/git-town/v17/pkg/set"
	"github.com/spf13/cobra"
)

const syncCommand = "sync"

const syncDesc = "Update the current branch with all relevant changes"

const syncHelp = `
Synchronizes the current branch with the rest of the world.

When run on a feature branch:
- syncs all ancestor branches
- pulls updates for the current branch
- merges the parent branch into the current branch
- pushes the current branch

When run on the main branch or a perennial branch:
- pulls and pushes updates for the current branch
- pushes tags

If the repository contains an "upstream" remote, syncs the main branch with its upstream counterpart. You can disable this by running "git config %s false".`

func Cmd() *cobra.Command {
	addAllFlag, readAllFlag := flags.All("sync all local branches")
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addNoPushFlag, readNoPushFlag := flags.NoPush()
	addStackFlag, readStackFlag := flags.Stack("sync the stack that the current branch belongs to")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     syncCommand,
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   syncDesc,
		Long:    cmdhelpers.Long(syncDesc, fmt.Sprintf(syncHelp, configdomain.KeySyncUpstream)),
		RunE: func(cmd *cobra.Command, _ []string) error {
			allBranches, err := readAllFlag(cmd)
			if err != nil {
				return err
			}
			detached, err := readDetachedFlag(cmd)
			if err != nil {
				return err
			}
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			noPush, err := readNoPushFlag(cmd)
			if err != nil {
				return err
			}
			stack, err := readStackFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeSync(allBranches, stack, detached, dryRun, verbose, noPush)
		},
	}
	addAllFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addNoPushFlag(&cmd)
	addStackFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSync(syncAllBranches configdomain.AllBranches, syncStack configdomain.FullStack, detached configdomain.Detached, dryRun configdomain.DryRun, verbose configdomain.Verbose, pushBranches configdomain.PushBranches) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineSyncData(syncAllBranches, syncStack, repo, verbose, detached)
	if err != nil || exit {
		return err
	}
	err = validateSyncData(data)
	if err != nil {
		return err
	}
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, repo.FinalMessages)
	runProgram := NewMutable(&program.Program{})
	branchesToDelete := set.New[gitdomain.LocalBranchName]()
	BranchesProgram(data.branchesToSync, BranchProgramArgs{
		BranchInfos:         data.branchInfos,
		BranchesToDelete:    NewMutable(&branchesToDelete),
		Config:              data.config,
		InitialBranch:       data.initialBranch,
		PrefetchBranchInfos: data.prefetchBranchesSnapshot.Branches,
		Program:             runProgram,
		PushBranches:        pushBranches,
		Remotes:             data.remotes,
	})
	previousbranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	finalBranchCandidates := gitdomain.LocalBranchNames{data.initialBranch}
	if previousBranch, hasPreviousBranch := data.previousBranch.Get(); hasPreviousBranch {
		finalBranchCandidates = append(finalBranchCandidates, previousBranch)
	}
	runProgram.Value.Add(&opcodes.CheckoutFirstExisting{
		Branches:   finalBranchCandidates,
		MainBranch: data.config.ValidatedConfigData.MainBranch,
	})
	if data.remotes.HasRemote(data.config.NormalConfig.DevRemote) && data.shouldPushTags && data.config.NormalConfig.IsOnline() {
		runProgram.Value.Add(&opcodes.PushTags{})
	}
	cmdhelpers.Wrap(runProgram, cmdhelpers.WrapOptions{
		DryRun:                   dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousbranchCandidates,
	})
	optimizedProgram := optimizer.Optimize(runProgram.Immutable())
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        0,
		Command:               syncCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            optimizedProgram,
		TouchedBranches:       optimizedProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               None[hostingdomain.Connector](),
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type syncData struct {
	branchInfos              gitdomain.BranchInfos
	branchesSnapshot         gitdomain.BranchesSnapshot
	branchesToSync           configdomain.BranchesToSync
	config                   config.ValidatedConfig
	detached                 configdomain.Detached
	dialogTestInputs         components.TestInputs
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	nonExistingBranches      gitdomain.LocalBranchNames
	prefetchBranchesSnapshot gitdomain.BranchesSnapshot
	previousBranch           Option[gitdomain.LocalBranchName]
	remotes                  gitdomain.Remotes
	shouldPushTags           bool
	stashSize                gitdomain.StashSize
}

func determineSyncData(syncAllBranches configdomain.AllBranches, syncStack configdomain.FullStack, repo execute.OpenRepoResult, verbose configdomain.Verbose, detached configdomain.Detached) (data syncData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	preFetchBranchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	previousBranch, hasPreviousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend).Get()
	var previousBranchOpt Option[gitdomain.LocalBranchName]
	if hasPreviousBranch {
		if previousBranchInfo, hasPreviousBranchInfo := branchesSnapshot.Branches.FindByLocalName(previousBranch).Get(); hasPreviousBranchInfo {
			switch previousBranchInfo.SyncStatus {
			case
				gitdomain.SyncStatusLocalOnly,
				gitdomain.SyncStatusNotInSync,
				gitdomain.SyncStatusAhead,
				gitdomain.SyncStatusBehind,
				gitdomain.SyncStatusUpToDate:
				previousBranchOpt = previousBranchInfo.LocalName
			case
				gitdomain.SyncStatusDeletedAtRemote,
				gitdomain.SyncStatusRemoteOnly,
				gitdomain.SyncStatusOtherWorktree:
				previousBranchOpt = None[gitdomain.LocalBranchName]()
			}
		}
	} else {
		previousBranchOpt = None[gitdomain.LocalBranchName]()
	}
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	perennialAndMain := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	var branchNamesToSync gitdomain.LocalBranchNames
	switch {
	case syncAllBranches.Enabled():
		branchNamesToSync = localBranches
	case syncStack.Enabled():
		branchNamesToSync = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialAndMain)
	default:
		branchNamesToSync = gitdomain.LocalBranchNames{initialBranch}
	}
	branchesAndTypes = repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err = validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: branchNamesToSync,
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	var shouldPushTags bool
	switch {
	case validatedConfig.NormalConfig.SyncTags.IsFalse():
		shouldPushTags = false
	case syncAllBranches.Enabled():
		shouldPushTags = true
	default:
		shouldPushTags = validatedConfig.IsMainOrPerennialBranch(initialBranch)
	}
	allBranchNamesToSync := validatedConfig.NormalConfig.Lineage.BranchesAndAncestors(branchNamesToSync)
	if detached {
		allBranchNamesToSync = validatedConfig.RemovePerennials(allBranchNamesToSync)
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, allBranchNamesToSync...)
	branchesToSync, err := BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, false, err
	}
	return syncData{
		branchInfos:              branchesSnapshot.Branches,
		branchesSnapshot:         branchesSnapshot,
		branchesToSync:           branchesToSync,
		config:                   validatedConfig,
		detached:                 detached,
		dialogTestInputs:         dialogTestInputs,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		nonExistingBranches:      nonExistingBranches,
		prefetchBranchesSnapshot: preFetchBranchesSnapshot,
		previousBranch:           previousBranchOpt,
		remotes:                  remotes,
		shouldPushTags:           shouldPushTags,
		stashSize:                stashSize,
	}, false, err
}

func BranchesToSync(branchInfosToSync gitdomain.BranchInfos, allBranchInfos gitdomain.BranchInfos, repo execute.OpenRepoResult, mainBranch gitdomain.LocalBranchName) (configdomain.BranchesToSync, error) {
	result := make(configdomain.BranchesToSync, len(branchInfosToSync))
	for b, branchInfo := range branchInfosToSync {
		branchNameToSync := branchInfo.GetLocalOrRemoteName()
		if branchNameToSync.LocalName() == mainBranch {
			result[b] = configdomain.BranchToSync{
				BranchInfo:         branchInfo,
				FirstCommitMessage: None[gitdomain.CommitMessage](),
			}
			continue
		}
		parentLocalName, hasParentName := repo.UnvalidatedConfig.NormalConfig.Lineage.Parent(branchNameToSync.LocalName()).Get()
		if !hasParentName {
			parentLocalName = mainBranch
		}
		parentBranchInfo, hasParentBranchInfo := allBranchInfos.FindLocalOrRemote(parentLocalName, repo.UnvalidatedConfig.NormalConfig.DevRemote).Get()
		if !hasParentBranchInfo {
			result[b] = configdomain.BranchToSync{
				BranchInfo:         branchInfo,
				FirstCommitMessage: None[gitdomain.CommitMessage](),
			}
			continue
		}
		parentBranchName := parentBranchInfo.GetLocalOrRemoteName()
		firstCommitMessage, err := repo.Git.FirstCommitMessageInBranch(repo.Backend, branchNameToSync, parentBranchName)
		if err != nil {
			return result, err
		}
		result[b] = configdomain.BranchToSync{
			BranchInfo:         branchInfo,
			FirstCommitMessage: firstCommitMessage,
		}
	}
	return result, nil
}

func validateSyncData(data syncData) error {
	// ensure any branch that uses the ff-only sync strategy does not have unpushed local commits
	if data.config.NormalConfig.SyncPerennialStrategy == configdomain.SyncPerennialStrategyFFOnly {
		perennialBranchesToSync := data.config.BranchesOfType(data.branchesToSync.BranchNames(), configdomain.BranchTypePerennialBranch)
		for _, perennialBranchToSync := range perennialBranchesToSync {
			if branchInfo, hasBranchInfo := data.branchInfos.FindByLocalName(perennialBranchToSync).Get(); hasBranchInfo {
				switch branchInfo.SyncStatus {
				case gitdomain.SyncStatusAhead, gitdomain.SyncStatusNotInSync:
					return fmt.Errorf(messages.SyncPerennialBranchHasUnpushedCommits, perennialBranchToSync)
				case gitdomain.SyncStatusBehind, gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusLocalOnly, gitdomain.SyncStatusOtherWorktree, gitdomain.SyncStatusRemoteOnly, gitdomain.SyncStatusUpToDate:
					// no problem with these sync statuses
				}
			}
		}
	}
	return nil
}
