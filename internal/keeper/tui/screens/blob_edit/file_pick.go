package blobedit

import (
	"fmt"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/screens"
	"gophkeeper/internal/keeper/tui/styles"
	"gophkeeper/pkg/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type FilePickScreen struct {
	secret  *models.Secret
	storage storage.Storage

	filePicker filepicker.Model
	callback   tui.NavigationCallback
}

func (s FilePickScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewFilePickScreen(msg.Secret, msg.Storage, msg.Callback), nil
}

func NewFilePickScreen(secret *models.Secret, strg storage.Storage, callback tui.NavigationCallback) *FilePickScreen {
	defaultPath, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting working directory: %v\n")
	}

	fp := filepicker.New()
	fp.CurrentDirectory = filepath.Join(defaultPath)
	fp.AutoHeight = false
	fp.Height = 10

	m := FilePickScreen{
		filePicker: fp,
		secret:     secret,
		storage:    strg,
		callback:   callback,
	}

	return &m
}

func (s FilePickScreen) Init() tea.Cmd {
	return s.filePicker.Init()
}

func (s *FilePickScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			cmds = append(cmds, tui.SetBodyPane(tui.BlobEditScreen, tui.WithStorage(s.storage), tui.WithSecret(s.secret)))
		}
	case tea.WindowSizeMsg:
		s.filePicker.Height = msg.Height - styles.FilepickerBotPadding
	}

	s.filePicker, cmd = s.filePicker.Update(msg)
	cmds = append(cmds, cmd)

	if selected, path := s.filePicker.DidSelectFile(msg); selected {
		cmds = append(cmds, tui.ReportInfo("selected: %v", path))
		cmds = append(cmds, s.callback(path))
	}

	return tea.Batch(cmds...)
}

func (s FilePickScreen) View() string {

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%20s%s:\n", "", s.filePicker.CurrentDirectory))
	b.WriteString(s.filePicker.View())

	return screens.RenderContent("Select file to store. Use ←, ↑, →, ↓ to navigate. Press B to go back", b.String())
}
