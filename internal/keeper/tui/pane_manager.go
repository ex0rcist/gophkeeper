package tui

import (
	"errors"
	"gophkeeper/internal/keeper/tui/styles"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/maps"
)

type Screen int

const (
	MenuScreen Screen = iota
	WelcomeScreen
	StorageCreateScreen
	StorageOpenScreen
	StorageBrowseScreen
	SecretTypeScreen
	FilePickScreen

	CredentialEditScreen
	TextEditScreen
	CardEditScreen
	BlobEditScreen
)

var (
	NoMakerError = errors.New("No maker for requested screen")
)

type Position int

const (
	BodyPane Position = iota
	LeftPane
)

const borderSize = 2

type PaneManager struct {
	makers        map[Screen]ScreenMaker
	cache         *Cache            // cache of previously made models   ????????
	focused       Position          // the position of the currently focused pane
	panes         map[Position]pane // panes tracks currently visible panes
	width, height int               // total width and height of the terminal space available to panes.
}

type pane struct {
	model Teable
	page  Page
}

func NewPaneManager(makers map[Screen]ScreenMaker) *PaneManager {
	p := &PaneManager{
		makers: makers,
		cache:  NewCache(),
		panes:  make(map[Position]pane),
	}
	return p
}

func (p *PaneManager) Init() tea.Cmd {
	return tea.Batch(
		SetLeftPane(MenuScreen),
		SetBodyPane(WelcomeScreen),
	)
}

func (p *PaneManager) Update(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, GlobalKeys.Tab):
			p.cycleFocusedPane()
		default:
			// Send remaining keys to focused pane
			cmds = append(cmds, p.updateModel(p.focused, msg))
		}

	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height

		p.updateChildSizes()
	case NavigationMsg:
		cmds = append(cmds, p.setPane(msg))
	default:
		// Send remaining message types to cached panes. // ?????
		cmds = p.cache.UpdateAll(msg)
	}

	return tea.Batch(cmds...)
}

func (pm *PaneManager) FocusedModel() Teable {
	return pm.panes[pm.focused].model
}

func (pm *PaneManager) cycleFocusedPane() {
	positions := maps.Keys(pm.panes)
	slices.Sort(positions)

	focusedIndex := int(pm.focused)
	totalPanes := len(pm.panes)

	if focusedIndex >= totalPanes-1 {
		focusedIndex = 0
	} else {
		focusedIndex++
	}

	pm.focusPane(positions[focusedIndex])
}

func (pm *PaneManager) updateChildSizes() {
	for position := range pm.panes {
		pm.updateModel(position, tea.WindowSizeMsg{
			Width:  pm.paneWidth(position) - borderSize,
			Height: pm.paneHeight(position) - borderSize,
		})
	}
}

func (pm *PaneManager) updateModel(position Position, msg tea.Msg) tea.Cmd {
	if pane, ok := pm.panes[position]; ok {
		return pane.model.Update(msg)
	}

	return nil
}

func (pm *PaneManager) setPane(msg NavigationMsg) tea.Cmd {
	var (
		cmd tea.Cmd
	)

	if pane, ok := pm.panes[msg.Position]; ok && pane.page == msg.Page {
		// Pane is already showing requested page, so just bring it into focus.
		if !msg.DisableFocus {
			pm.focusPane(msg.Position)
		}

		return nil
	}

	model := pm.cache.Get(msg.Page)
	if 1 == 1 || model == nil { // TODO!!!

		maker, ok := pm.makers[msg.Page.Screen]
		if !ok {
			return ReportError(NoMakerError)
		}

		var err error
		model, err = maker.Make(msg, 0, 0)
		if err != nil {
			return ReportError(err)
		}

		pm.cache.Put(msg.Page, model)
		cmd = model.Init()
	}

	pm.panes[msg.Position] = pane{model: model, page: msg.Page}
	pm.updateChildSizes()

	if !msg.DisableFocus {
		pm.focusPane(msg.Position)
	}

	return cmd
}

func (pm *PaneManager) focusPane(position Position) {
	if _, ok := pm.panes[position]; ok {
		pm.focused = position
	}
}

func (pm *PaneManager) paneWidth(position Position) int {
	switch position {
	case LeftPane:
		return defaultLeftPaneWidth
	case BodyPane:
		return pm.width - pm.paneWidth(LeftPane)
	default:
		return pm.width
	}
}

func (pm *PaneManager) paneHeight(_ Position) int {
	return pm.height
}

func (pm *PaneManager) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top,
			pm.renderPane(LeftPane),
			pm.renderPane(BodyPane),
		),
	)
}

func (m *PaneManager) renderPane(position Position) string {
	if _, ok := m.panes[position]; !ok {
		return ""
	}

	// Width and Height does not include border size, so substract it
	paneStyle := styles.InactiveBorder.
		Width(m.paneWidth(position) - borderSize).
		Height(m.paneHeight(position) - borderSize)

	if position == m.focused {
		paneStyle = styles.ActiveBorder.Inherit(paneStyle)
	}

	model := m.panes[position].model
	return paneStyle.Render(model.View())

}

func (pm *PaneManager) HelpBindings() (bindings []key.Binding) {
	if model, ok := pm.FocusedModel().(ModelHelpBindings); ok {
		bindings = append(bindings, model.HelpBindings()...)
	}
	return bindings
}
