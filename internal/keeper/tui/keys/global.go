package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type global struct {
	Logs key.Binding
	Quit key.Binding
	Help key.Binding
	Tab  key.Binding
	Back key.Binding
}

var Global = global{
	Logs: key.NewBinding(
		key.WithKeys("ctrl+l", "fn+l"),
		key.WithHelp("ctrl+l", "show logs"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"), // TODO: intercept global "?"
		key.WithHelp("ctrl+h", "help"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "cycle panes"),
	),
	Back: key.NewBinding(
		key.WithKeys("ctrl+b"),
		key.WithHelp("ctrl+b", "back"),
	),
}
