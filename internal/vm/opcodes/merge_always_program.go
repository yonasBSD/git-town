package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	"github.com/git-town/git-town/v19/internal/vm/shared"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

// MergeAlwaysProgram merges the feature branch into its parent by always creating a merge comment (merge --no-ff).
type MergeAlwaysProgram struct {
	Branch        gitdomain.LocalBranchName
	CommitMessage Option[gitdomain.CommitMessage]
	undeclaredOpcodeMethods
}

func (self *MergeAlwaysProgram) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *MergeAlwaysProgram) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *MergeAlwaysProgram) Run(args shared.RunArgs) error {
	// Reverting parent is intentionally not supported due to potential confusion
	// caused by reverted merge commit. See
	// <https://github.com/git/git/blob/master/Documentation/howto/revert-a-faulty-merge.txt>
	// for more information.
	useMessage := configdomain.UseCustomMessageOr(self.CommitMessage, configdomain.EditDefaultMessage())
	return args.Git.MergeNoFastForward(args.Frontend, useMessage, self.Branch)
}

func (self *MergeAlwaysProgram) ShouldUndoOnError() bool {
	return true
}
