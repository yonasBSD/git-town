package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	syncPrototypeStrategyTitle = `Sync-prototype strategy`
	SyncPrototypeStrategyHelp  = `
Choose how Git Town should
synchronize prototype branches.

Prototype branches are local-only feature branches.
They are useful for reducing load on CI systems
and limiting the sharing of confidential changes.

`
)

func SyncPrototypeStrategy(existing configdomain.SyncPrototypeStrategy, inputs components.TestInput) (configdomain.SyncPrototypeStrategy, dialogdomain.Exit, error) {
	entries := list.Entries[configdomain.SyncPrototypeStrategy]{
		{
			Data: configdomain.SyncPrototypeStrategyMerge,
			Text: "merge updates from the parent and tracking branch",
		},
		{
			Data: configdomain.SyncPrototypeStrategyRebase,
			Text: "rebase branches against their parent and tracking branch",
		},
		{
			Data: configdomain.SyncPrototypeStrategyCompress,
			Text: "compress the branch after merging parent and tracking",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, exit, err := components.RadioList(entries, defaultPos, syncPrototypeStrategyTitle, SyncPrototypeStrategyHelp, inputs)
	if err != nil || exit {
		return configdomain.SyncPrototypeStrategyMerge, exit, err
	}
	fmt.Printf(messages.SyncPrototypeBranches, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
