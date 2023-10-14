package step

import "github.com/git-town/git-town/v9/src/domain"

// RebaseParent rebases the given branch against the branch that is its parent at runtime.
type RebaseParent struct {
	CurrentBranch domain.LocalBranchName
	Empty
}

func (step *RebaseParent) CreateAbortSteps() []Step {
	return []Step{
		&AbortRebase{},
	}
}

func (step *RebaseParent) CreateContinueSteps() []Step {
	return []Step{
		&ContinueRebase{},
	}
}

func (step *RebaseParent) Run(args RunArgs) error {
	parent := args.Lineage.Parent(step.CurrentBranch)
	if parent.IsEmpty() {
		return nil
	}
	return args.Runner.Frontend.Rebase(parent.BranchName())
}
