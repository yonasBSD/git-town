package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// SetParentIfBranchExists sets the given parent branch as the parent of the given branch,
// but only the latter exists.
type SetParentIfBranchExists struct {
	Branch                  gitdomain.LocalBranchName
	Parent                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SetParentIfBranchExists) Run(args shared.RunArgs) error {
	if !args.Git.BranchExists(args.Backend, self.Branch) {
		return nil
	}
	return args.Config.SetParent(self.Branch, self.Parent)
}
