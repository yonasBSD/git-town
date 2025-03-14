package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/messages"
)

const (
	syncPerennialStrategyTitle = `Sync-perennial strategy`
	SyncPerennialStrategyHelp  = `
How should Git Town synchronize perennial branches?
Perennial branches have no parent branch.
The only updates they receive are additional commits
made to their tracking branch somewhere else.

`
)

func SyncPerennialStrategy(existing configdomain.SyncPerennialStrategy, inputs components.TestInput) (configdomain.SyncPerennialStrategy, bool, error) {
	entries := list.Entries[configdomain.SyncPerennialStrategy]{
		{
			Data: configdomain.SyncPerennialStrategyFFOnly,
			Text: "fast-forward perennial branches to their tracking branch",
		},
		{
			Data: configdomain.SyncPerennialStrategyRebase,
			Text: "rebase perennial branches against their tracking branch",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, syncPerennialStrategyTitle, SyncPerennialStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.SyncPerennialStrategyRebase, aborted, err
	}
	fmt.Printf(messages.SyncPerennialBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
