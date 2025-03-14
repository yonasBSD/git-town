package shared

import (
	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/config"
	"github.com/git-town/git-town/v18/internal/forge/forgedomain"
	"github.com/git-town/git-town/v18/internal/git"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

type RunArgs struct {
	Backend                         gitdomain.RunnerQuerier
	BranchInfos                     Option[gitdomain.BranchInfos]
	Config                          Mutable[config.ValidatedConfig]
	Connector                       Option[forgedomain.Connector]
	DialogTestInputs                components.TestInputs
	FinalMessages                   stringslice.Collector
	Frontend                        gitdomain.Runner
	Git                             git.Commands
	PrependOpcodes                  func(...Opcode)
	RegisterUndoablePerennialCommit func(gitdomain.SHA)
	UpdateInitialSnapshotLocalSHA   func(gitdomain.LocalBranchName, gitdomain.SHA) error
}
