package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Teable interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	View() string
}

type ScreenMaker interface {
	// msg.Page.ID ?
	Make(msg NavigationMsg, width, height int) (Teable, error)
}

// Page identifies an instance of a model
type Page struct {
	// The model kind. Identifies the model maker to construct the page.
	Screen Screen
	// The ID of the resource for a model. In the case of global listings of
	// modules, workspaces, etc, this is the global resource.
	ID string //resource.ID
}

// ModelHelpBindings is implemented by models that surface further help bindings
// specific to the model.
type ModelHelpBindings interface {
	HelpBindings() []key.Binding
}
