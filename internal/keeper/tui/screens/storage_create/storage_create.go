package storagecreate

import (
	"fmt"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/usecase"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StorageCreateScreen struct {
	style lipgloss.Style

	textInput textinput.Model
}

func (s StorageCreateScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewStorageCreateScreen(), nil
}

func NewStorageCreateScreen() *StorageCreateScreen {
	ti := textinput.New()
	ti.Placeholder = "default.txt"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &StorageCreateScreen{
		textInput: ti,
		style:     lipgloss.NewStyle(),
	}
}

func (s StorageCreateScreen) Init() tea.Cmd {
	return textinput.Blink

}

func (s *StorageCreateScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			path := s.textInput.Value()

			uc := usecase.NewCreateStorageUsecase()
			_, err := uc.Call(path)
			if err != nil {
				cmds = append(cmds, tui.ReportInfo("Error: %v", err))
			} else {
				cmds = append(cmds, tui.ReportInfo("Created new storage: %v", path))
			}

			// TODO: clear history? go to editing storage screen
		}

		s.textInput, cmd = s.textInput.Update(msg)
	}

	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (s StorageCreateScreen) View() string {
	return fmt.Sprintf(
		"Enter file name:\n\n%s\n",
		s.textInput.View(),
	)

}
