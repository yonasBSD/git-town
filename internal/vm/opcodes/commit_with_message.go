package opcodes

import (
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/vm/shared"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// CommitWithMessage commits all open changes using the given commit message.
type CommitWithMessage struct {
	AuthorOverride          Option[gitdomain.Author]
	Message                 gitdomain.CommitMessage
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitWithMessage) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, configdomain.UseCustomMessage(self.Message), self.AuthorOverride)
}
