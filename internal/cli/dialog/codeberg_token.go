package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	codebergTokenTitle = `Codeberg API token`
	codebergTokenHelp  = `
Git Town can update pull requests
and ship branches on codeberg-based forges for you.
To enable this, please enter a codeberg API token.
More info at
https://docs.codeberg.org/advanced/access-token.

If you leave this empty,
Git Town will not use the codeberg API.

`
)

// CodebergToken lets the user enter the Gitea API token.
func CodebergToken(oldValue Option[forgedomain.CodebergToken], inputs components.TestInput) (Option[forgedomain.CodebergToken], dialogdomain.Exit, error) {
	text, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          codebergTokenHelp,
		Prompt:        "Your Codeberg API token: ",
		TestInput:     inputs,
		Title:         codebergTokenTitle,
	})
	fmt.Printf(messages.CodebergToken, components.FormattedSecret(text, aborted))
	return forgedomain.ParseCodebergToken(text), aborted, err
}
