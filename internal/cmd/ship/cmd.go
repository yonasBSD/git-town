package ship

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/flags"
	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/execute"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/gohacks/stringslice"
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
	"github.com/spf13/cobra"
)

const shipCommand = "ship"

const shipDesc = "Deliver a completed feature branch"

const shipHelp = `
Merges the given or current feature branch into its parent.
How exactly this happen depends on the configured ship-strategy.

Ships only direct children of the main branch.
To ship a child branch, ship or delete all ancestor branches first
or ship with the "--to-parent" flag.

To use the online functionality, configure a personal access token with the "repo" scope
and run "git config %s <token>" (optionally add the "--global" flag).

If your origin server deletes shipped branches,
disable the ship-delete-tracking-branch configuration setting.`

func Cmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addMessageFlag, readMessageFlag := flags.CommitMessage("specify the commit message for the squash commit")
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addShipStrategyFlag, readShipStrategyFlag := flags.ShipStrategy()
	addToParentFlag, readToParentFlag := flags.ShipIntoNonPerennialParent()
	cmd := cobra.Command{
		Use:   shipCommand,
		Args:  cobra.MaximumNArgs(1),
		Short: shipDesc,
		Long:  cmdhelpers.Long(shipDesc, fmt.Sprintf(shipHelp, configdomain.KeyGithubToken)),
		RunE: func(cmd *cobra.Command, args []string) error {
			shipStrategyOverride, err := readShipStrategyFlag(cmd)
			if err != nil {
				return err
			}
			message, err := readMessageFlag(cmd)
			if err != nil {
				return err
			}
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			toParent, err := readToParentFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeShip(args, message, dryRun, verbose, shipStrategyOverride, toParent)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addMessageFlag(&cmd)
	addShipStrategyFlag(&cmd)
	addToParentFlag(&cmd)
	return &cmd
}

func executeShip(args []string, message Option[gitdomain.CommitMessage], dryRun configdomain.DryRun, verbose configdomain.Verbose, shipStrategy Option[configdomain.ShipStrategy], toParent configdomain.ShipIntoNonperennialParent) error {
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
	sharedData, exit, err := determineSharedShipData(args, repo, dryRun, shipStrategy, verbose)
	if err != nil || exit {
		return err
	}
	err = validateSharedData(sharedData, toParent, message)
	if err != nil {
		return err
	}
	prog := NewMutable(&program.Program{})
	switch sharedData.config.NormalConfig.ShipStrategy {
	case configdomain.ShipStrategyAPI:
		apiData, err := determineAPIData(sharedData)
		if err != nil {
			return err
		}
		err = shipAPIProgram(prog, sharedData, apiData, message)
		if err != nil {
			return err
		}
	case configdomain.ShipStrategyAlwaysMerge:
		mergeData, err := determineMergeData(repo, sharedData.branchNameToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramAlwaysMerge(prog, sharedData, mergeData, message)
	case configdomain.ShipStrategyFastForward:
		mergeData, err := determineMergeData(repo, sharedData.branchNameToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramFastForward(prog, sharedData, mergeData)
	case configdomain.ShipStrategySquashMerge:
		squashMergeData, err := determineMergeData(repo, sharedData.branchNameToShip, sharedData.targetBranchName)
		if err != nil {
			return err
		}
		shipProgramSquashMerge(prog, sharedData, squashMergeData, message)
	}
	optimizedProgram := optimizer.Optimize(prog.Immutable())
	runState := runstate.RunState{
		BeginBranchesSnapshot: sharedData.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        sharedData.stashSize,
		Command:               shipCommand,
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
		Config:                  sharedData.config,
		Connector:               sharedData.connector,
		DialogTestInputs:        sharedData.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          sharedData.hasOpenChanges,
		InitialBranch:           sharedData.initialBranch,
		InitialBranchesSnapshot: sharedData.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        sharedData.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

func UpdateChildBranchProposalsToGrandParent(prog *program.Program, proposals []hostingdomain.Proposal) {
	for _, childProposal := range proposals {
		prog.Add(&opcodes.ProposalUpdateTargetToGrandParent{
			Branch:         childProposal.Source,
			OldTarget:      childProposal.Target,
			ProposalNumber: childProposal.Number,
		})
	}
}

func validateSharedData(data sharedShipData, toParent configdomain.ShipIntoNonperennialParent, message Option[gitdomain.CommitMessage]) error {
	if data.config.NormalConfig.ShipStrategy == configdomain.ShipStrategyFastForward && message.IsSome() {
		return errors.New(messages.ShipMessageWithFastForward)
	}
	if !toParent {
		branch := data.branchToShip.LocalName.GetOrPanic()
		parentBranch := data.targetBranch.LocalName.GetOrPanic()
		if !data.config.IsMainOrPerennialBranch(parentBranch) {
			ancestors := data.config.NormalConfig.Lineage.Ancestors(branch)
			ancestorsWithoutMainOrPerennial := ancestors[1:]
			oldestAncestor := ancestorsWithoutMainOrPerennial[0]
			return fmt.Errorf(messages.ShipChildBranch, stringslice.Connect(ancestorsWithoutMainOrPerennial.Strings()), oldestAncestor)
		}
	}
	switch data.branchToShip.SyncStatus {
	case gitdomain.SyncStatusDeletedAtRemote:
		return fmt.Errorf(messages.BranchDeletedAtRemote, data.branchNameToShip)
	case
		gitdomain.SyncStatusNotInSync,
		gitdomain.SyncStatusAhead,
		gitdomain.SyncStatusBehind:
		return fmt.Errorf(messages.ShipBranchNotInSync, data.branchNameToShip)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.ShipBranchIsInOtherWorktree, data.branchNameToShip)
	case
		gitdomain.SyncStatusUpToDate,
		gitdomain.SyncStatusRemoteOnly,
		gitdomain.SyncStatusLocalOnly:
	}
	if localName, hasLocalName := data.branchToShip.LocalName.Get(); hasLocalName {
		if localName == data.initialBranch {
			return validate.NoOpenChanges(data.hasOpenChanges)
		}
	}
	return nil
}