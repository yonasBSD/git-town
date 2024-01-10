package execute

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/validate"
)

// EnsureKnownBranchAncestry makes sure the lineage for the given branch is known.
// If needed, it queries the user for missing information.
// It returns the updated version of all information that is derived from the lineage.
//
// The purpose of this function is to implement proper cache invalidation.
// It ensures that all information derived from lineage gets updated when the lineage is updated.
func EnsureKnownBranchAncestry(branch gitdomain.LocalBranchName, args EnsureKnownBranchAncestryArgs) error {
	updated, err := validate.KnowsBranchAncestors(branch, validate.KnowsBranchAncestorsArgs{
		AllBranches: args.AllBranches.Names(),
		Backend:     &args.Runner.Backend,
		Config:      args.Config,
		MainBranch:  args.DefaultBranch,
	})
	if err != nil {
		return err
	}
	if updated {
		// reload after ancestry change
		args.Runner.Config.Reload()
	}
	return nil
}

type EnsureKnownBranchAncestryArgs struct {
	Config        *configdomain.FullConfig
	AllBranches   gitdomain.BranchInfos
	DefaultBranch gitdomain.LocalBranchName
	Runner        *git.ProdRunner
}
