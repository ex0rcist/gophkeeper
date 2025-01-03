package welcome

import (
	"gophkeeper/internal/keeper/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeScreen struct {
	style lipgloss.Style
}

func (w WelcomeScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewWelcomeScreen(), nil
}

func NewWelcomeScreen() *WelcomeScreen {
	return &WelcomeScreen{
		style: lipgloss.NewStyle(),
	}
}

func (m WelcomeScreen) Init() tea.Cmd {
	return tea.SetWindowTitle("GophKeeper client")
}

func (m *WelcomeScreen) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (m WelcomeScreen) View() string {

	c := `
   _____             _     _  __                         
  / ____|           | |   | |/ /                         
 | |  __  ___  _ __ | |__ | ' / ___  ___ _ __   ___ _ __ 
 | | |_ |/ _ \| '_ \| '_ \|  < / _ \/ _ \ '_ \ / _ \ '__|
 | |__| | (_) | |_) | | | | . \  __/  __/ |_) |  __/ |   
  \_____|\___/| .__/|_| |_|_|\_\___|\___| .__/ \___|_|   
              | |                       | |              
              |_|                       |_|              
	`

	return m.style.Render(c)
}
