package dialog

import "github.com/muesli/termenv"

// Typical colors used in BubbleTea dialogs.
type dialogColors struct {
	Help      termenv.Style // color of help text
	HelpKey   termenv.Style // color of key names in help text
	Initial   termenv.Style // color for the row containing the currently checked out branch
	Selection termenv.Style // color for the currently selected entry
}

func FormattedSelection(selection string, aborted bool) string {
	if aborted {
		return red().Styled("(aborted)")
	}
	return green().Styled(selection)
}

func createColors() dialogColors {
	return dialogColors{
		Help:      termenv.String().Faint(),
		HelpKey:   termenv.String().Faint().Bold(),
		Initial:   termenv.String().Foreground(termenv.ANSIGreen),
		Selection: termenv.String().Foreground(termenv.ANSICyan),
	}
}

func green() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIGreen)
}

func red() termenv.Style {
	return termenv.String().Foreground(termenv.ANSIRed)
}
