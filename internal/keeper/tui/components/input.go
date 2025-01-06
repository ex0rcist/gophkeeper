package components

import (
	"gophkeeper/internal/keeper/tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
)

type inputOpts struct {
	placeholder string
	charLimit   int
	focus       bool
}

func NewInput(opts inputOpts) textinput.Model {
	t := textinput.New()
	t.CharLimit = opts.charLimit
	t.Placeholder = opts.placeholder

	if opts.focus {
		t.Focus()
		t.PromptStyle = styles.Focused
		t.TextStyle = styles.Focused
	}

	return t
}
