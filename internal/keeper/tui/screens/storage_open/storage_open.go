package storageopen

import (
	"fmt"
	"gophkeeper/internal/keeper/crypto"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/styles"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type passwordProvidedMsg struct {
	path     string
	password string
}

type StorageOpenScreen struct {
	filePicker filepicker.Model
	encrypter  crypto.Encrypter
	selected   string
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
		encrypter:  crypto.NewKeeperEncrypter(),
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

	switch msg := msg.(type) {
	case passwordProvidedMsg: // msg from prompt for password
		strg, err := storage.NewFileStorage(msg.path, msg.password, s.encrypter)
		if err != nil {
			cmds = append(cmds, tui.ReportError(err))
		} else {
			cmd = tui.NavigateTo(tui.StorageBrowseScreen, tui.WithStorage(strg), tui.WithPosition(tui.BodyPane))
			cmds = append(cmds, cmd)
		}
	default:
		s.filePicker, cmd = s.filePicker.Update(msg)
		cmds = append(cmds, cmd)

		if selected, path := s.filePicker.DidSelectFile(msg); selected {
			return tui.StringPrompt("enter password", func(str string) tea.Cmd {
				return func() tea.Msg { return passwordProvidedMsg{path: path, password: str} }
			})
		}
	}

	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (s StorageOpenScreen) View() string {
	var b strings.Builder

	titleStyle := styles.Bold.Foreground(lipgloss.Color("#FF79C6"))

	b.WriteString(titleStyle.Render("Select storage file to open. Use ←, ↑, →, ↓ to navigate"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("%20s%s:\n", "", s.filePicker.CurrentDirectory))
	b.WriteString(s.filePicker.View())

	return b.String()
}
