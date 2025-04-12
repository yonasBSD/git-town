package format

import "github.com/git-town/git-town/v18/internal/git/gitdomain"

func BranchNames(branches gitdomain.LocalBranchNames) string {
	if len(branches) == 0 {
		return "(none)"
	}
	return branches.Join(", ")
}
