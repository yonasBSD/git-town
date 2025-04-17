package debug

import (
	"os"

	"github.com/git-town/git-town/v19/internal/cli/dialog"
	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterDefaultBranch() *cobra.Command {
	return &cobra.Command{
		Use: "default-branch-type",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.DefaultBranchType(configdomain.BranchTypeFeatureBranch, dialogTestInputs.Next())
			return err
		},
	}
}
