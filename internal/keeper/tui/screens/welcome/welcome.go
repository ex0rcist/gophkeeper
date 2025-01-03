package welcome

import (
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	logo = `
   _____             _     _  __                         
  / ____|           | |   | |/ /                         
 | |  __  ___  _ __ | |__ | ' / ___  ___ _ __   ___ _ __ 
 | | |_ |/ _ \| '_ \| '_ \|  < / _ \/ _ \ '_ \ / _ \ '__|
 | |__| | (_) | |_) | | | | . \  __/  __/ |_) |  __/ |   
  \_____|\___/| .__/|_| |_|_|\_\___|\___| .__/ \___|_|   
              | |                       | |              
              |_|                       |_|              
	`
)

type WelcomeScreen struct {
	width, height int
}

func (s WelcomeScreen) Make(_ tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewWelcomeScreen(width, height), nil
}

func NewWelcomeScreen(width, height int) *WelcomeScreen {
	return &WelcomeScreen{
		width:  width,
		height: height,
	}
}

func (s WelcomeScreen) Init() tea.Cmd {
	return tea.SetWindowTitle("GophKeeper client")
}

func (s *WelcomeScreen) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	}

	return nil
}

func (s WelcomeScreen) View() string {
	return lipgloss.Place(
		s.width, s.height,
		lipgloss.Center, lipgloss.Center,
		styles.Highlighted.Render(logo),
	)
}
