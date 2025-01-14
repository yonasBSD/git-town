package ship

import (
	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/opcodes"
	"github.com/git-town/git-town/v17/internal/vm/program"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

func shipProgramAlwaysMerge(prog Mutable[program.Program], sharedData sharedShipData, mergeData shipDataMerge, commitMessage Option[gitdomain.CommitMessage]) {
	prog.Value.Add(&opcodes.BranchEnsureShippableChanges{Branch: sharedData.branchNameToShip, Parent: sharedData.targetBranchName})
	if sharedData.initialBranch != sharedData.targetBranchName {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.targetBranchName})
	}
	if mergeData.remotes.HasRemote(sharedData.config.NormalConfig.DevRemote) && sharedData.config.NormalConfig.IsOnline() {
		UpdateChildBranchProposalsToGrandParent(prog.Value, sharedData.proposalsOfChildBranches)
	}
	prog.Value.Add(&opcodes.MergeAlwaysProgram{Branch: sharedData.branchNameToShip, CommitMessage: commitMessage})
	if mergeData.remotes.HasRemote(sharedData.config.NormalConfig.DevRemote) && sharedData.config.NormalConfig.IsOnline() {
		prog.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: sharedData.targetBranchName})
	}
	if !sharedData.dryRun {
		prog.Value.Add(&opcodes.LineageParentRemove{Branch: sharedData.branchNameToShip})
	}
	if branchToShipRemoteName, hasRemoteName := sharedData.branchToShip.RemoteName.Get(); hasRemoteName {
		if sharedData.config.NormalConfig.IsOnline() {
			if sharedData.config.NormalConfig.ShipDeleteTrackingBranch {
				prog.Value.Add(&opcodes.BranchTrackingDelete{Branch: branchToShipRemoteName})
			}
		}
	}
	for _, child := range sharedData.childBranches {
		prog.Value.Add(&opcodes.LineageParentSetToGrandParent{Branch: child})
	}
	if !sharedData.isShippingInitialBranch {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.initialBranch})
	}
	prog.Value.Add(&opcodes.BranchLocalDelete{Branch: sharedData.branchNameToShip})
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{sharedData.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   sharedData.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         !sharedData.isShippingInitialBranch && sharedData.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
}