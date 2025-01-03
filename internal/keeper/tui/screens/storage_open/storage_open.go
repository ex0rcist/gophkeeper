package storageopen

import (
	"fmt"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StorageOpenScreen struct {
	filePicker filepicker.Model
	selected   string
	style      lipgloss.Style
}

func (s StorageOpenScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewStorageOpenScreen(), nil
}

func NewStorageOpenScreen() *StorageOpenScreen {
	defaultPath, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting working directory: %v\n")
	}

	fp := filepicker.New()
	fp.CurrentDirectory = filepath.Join(defaultPath)

	return &StorageOpenScreen{
		filePicker: fp,
		style:      lipgloss.NewStyle(),
	}
}

func (s StorageOpenScreen) Init() tea.Cmd {
	return s.filePicker.Init()
}

func (s *StorageOpenScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	s.filePicker, cmd = s.filePicker.Update(msg)
	cmds = append(cmds, cmd)

	if selected, path := s.filePicker.DidSelectFile(msg); selected {
		cmds = append(cmds, tui.ReportInfo("selected: %v", path))

		strg, err := storage.NewFileStorage(path)
		if err != nil {
			cmds = append(cmds, tui.ReportError(err))
		} else {
			cmd = tui.NavigateTo(
				tui.StorageBrowseScreen,
				tui.WithStorage(strg),
				tui.WithPosition(tui.BodyPane),
			)

			cmds = append(cmds, cmd)
		}
	}

	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (s StorageOpenScreen) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF79C6"))

	b.WriteString(titleStyle.Render("Select storage file to open. Use ←, ↑, →, ↓ to navigate"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("%20s%s:\n", "", s.filePicker.CurrentDirectory))
	b.WriteString(s.filePicker.View())

	return b.String()
}
