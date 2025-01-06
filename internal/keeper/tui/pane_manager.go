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
	LoginScreen
	RegisterScreen
	RemoteOpenScreen

	CredentialEditScreen
	TextEditScreen
	CardEditScreen
	BlobEditScreen
)

var (
	ErrNoMaker = errors.New("no maker for requested screen")
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

func (pm *PaneManager) Init() tea.Cmd {
	return tea.Batch(
		SetLeftPane(MenuScreen),
		SetBodyPane(WelcomeScreen),
	)
}

func (pm *PaneManager) Update(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, GlobalKeys.Tab):
			pm.cycleFocusedPane()
		default:
			// Send remaining keys to focused pane
			cmds = append(cmds, pm.updateModel(pm.focused, msg))
		}

	case tea.WindowSizeMsg:
		pm.width = msg.Width
		pm.height = msg.Height

		pm.updateChildSizes()
	case NavigationMsg:
		cmds = append(cmds, pm.setPane(msg))
	default:
		// Send remaining message types to cached panes. // ?????
		cmds = pm.cache.UpdateAll(msg)
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

	if p, ok := pm.panes[position]; ok {
		return p.model.Update(msg)
	}

	return nil
}

func (pm *PaneManager) setPane(msg NavigationMsg) tea.Cmd {
	var (
		cmd tea.Cmd
	)

	if p, ok := pm.panes[msg.Position]; ok && p.page == msg.Page {
		// Pane is already showing requested page, so just bring it into focus.
		if !msg.DisableFocus {
			pm.focusPane(msg.Position)
		}

		return nil
	}

	// TODO: cache invalidation
	useCache := true

	model := pm.cache.Get(msg.Page)
	if useCache || model == nil {

		maker, ok := pm.makers[msg.Page.Screen]
		if !ok {
			return ReportError(ErrNoMaker)
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

func (pm *PaneManager) renderPane(position Position) string {
	if _, ok := pm.panes[position]; !ok {
		return ""
	}

	// Width and Height does not include border size, so substract it
	paneStyle := styles.InactiveBorder.
		Width(pm.paneWidth(position) - borderSize).
		Height(pm.paneHeight(position) - borderSize)

	if position == pm.focused {
		paneStyle = styles.ActiveBorder.Inherit(paneStyle)
	}

	model := pm.panes[position].model
	return paneStyle.Render(model.View())

}

func (pm *PaneManager) HelpBindings() (bindings []key.Binding) {
	if model, ok := pm.FocusedModel().(ModelHelpBindings); ok {
		bindings = append(bindings, model.HelpBindings()...)
	}
	return bindings
}
