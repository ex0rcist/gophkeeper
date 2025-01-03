package menu

import (
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	menuItems = []MenuItem{
		{name: "Welcome", cmd: tui.SetBodyPane(tui.WelcomeScreen)},
		{name: "Open passfile", cmd: tui.SetBodyPane(tui.StorageOpenScreen)},
		{name: "Create new passfile", cmd: tui.SetBodyPane(tui.StorageCreateScreen)},
	}
)

type MenuItem struct {
	name string
	cmd  tea.Cmd
}

type MenuScreen struct {
	tea.Model

	list list.Model

	choice    string
	menuItems []MenuItem
	itemsMap  map[string]int
}

func (s MenuScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewMenu(), nil
}

func NewMenu() *MenuScreen {
	m := MenuScreen{itemsMap: make(map[string]int, len(menuItems))}
	m.prepareMenuModel(menuItems)

	return &m
}

func (s MenuScreen) Init() tea.Cmd {
	return nil
}

func (s *MenuScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.list.SetWidth(msg.Width)
		s.list.SetHeight(msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := s.list.SelectedItem().(item)
			if ok {
				s.choice = string(i)
				idx := s.itemsMap[string(i)]
				cmds = append(cmds, menuItems[idx].cmd)
			}
		}
	}

	s.list, cmd = s.list.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (s MenuScreen) View() string {
	return styles.Regular.Render(s.list.View())
}

func (s *MenuScreen) prepareMenuModel(menuItems []MenuItem) {
	listItems := []list.Item{}
	for i, menuItem := range menuItems {
		listItems = append(listItems, item(menuItem.name))
		s.itemsMap[menuItem.name] = i // for ease of search
	}

	s.list = list.New(listItems, itemDelegate{}, 0, 0)

	s.list.SetShowStatusBar(false)
	s.list.SetFilteringEnabled(false)
	s.list.SetShowTitle(false)
	s.list.SetShowPagination(false)
	s.list.SetShowHelp(false)
	s.list.KeyMap.Quit.SetEnabled(false)
}
