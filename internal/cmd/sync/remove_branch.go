package sync

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/opcodes"
	"github.com/git-town/git-town/v17/internal/vm/program"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

func RemoveBranchConfiguration(args RemoveBranchConfigurationArgs) {
	childBranches := args.Lineage.Children(args.Branch)
	for _, child := range childBranches {
		args.Program.Value.Add(&opcodes.LineageParentSetToGrandParent{Branch: child})
	}
	args.Program.Value.Add(&opcodes.LineageParentRemove{Branch: args.Branch})
}

type RemoveBranchConfigurationArgs struct {
	Branch  gitdomain.LocalBranchName
	Lineage configdomain.Lineage
	Program Mutable[program.Program]
}
