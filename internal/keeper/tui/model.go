package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Any object that can act as tea.Model
type Teable interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	View() string
}

// Any constructor of screen in case of external dependencies,
// or screen itself in case of none dependencies
type ScreenMaker interface {
	Make(msg NavigationMsg, width, height int) (Teable, error)
}

// Page identifies an instance of a model
type Page struct {
	// Screen tyoe. Identifies the screen maker to construct the page.
	Screen Screen

	// TODO
	// The ID of the resource for a model. In the case of global listings of
	// modules, workspaces, etc, this is the global resource.
	// ID string
}

// ModelHelpBindings is implemented by models
// that pass up own help bindings specific to that model.
type ModelHelpBindings interface {
	HelpBindings() []key.Binding
}
