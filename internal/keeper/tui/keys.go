package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

var GlobalKeys = struct {
	Quit key.Binding
	Help key.Binding
	Tab  key.Binding
}{
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
		key.WithHelp("tab", "cycle menu/body"),
	),
}
