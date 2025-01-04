package storagecreate

import (
	"fmt"
	"gophkeeper/internal/keeper/crypto"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/components"
	"gophkeeper/internal/keeper/tui/screens"
	"gophkeeper/internal/keeper/tui/styles"
	"gophkeeper/internal/keeper/usecase"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	createPath = iota
	createPassword
)

type storageCreatedMsg struct{}

type StorageCreateScreen struct {
	inputGroup components.InputGroup

	pathInput     textinput.Model
	passwordInput textinput.Model

	storage   storage.Storage
	encrypter crypto.Encrypter

	createStorageUC *usecase.CreateLocalStoreUseCase
}

func (s StorageCreateScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewStorageCreateScreen()
}

func NewStorageCreateScreen() (*StorageCreateScreen, error) {
	var err error

	hdir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home dir: %s", err.Error())
	}

	scr := &StorageCreateScreen{
		encrypter:       crypto.NewKeeperEncrypter(), // todo: dependencies
		createStorageUC: usecase.NewCreateStorageUsecase(),
	}

	inputs := make([]textinput.Model, 2)
	inputs[createPath] = newInput(inputOpts{placeholder: "Path to store", charLimit: 256, value: fmt.Sprintf("%s/%s", hdir, "secret.db")})
	inputs[createPassword] = newInput(inputOpts{placeholder: "Password", charLimit: 64})

	buttons := []components.Button{}
	buttons = append(buttons, components.Button{Title: "[ Submit ]", Cmd: func() tea.Cmd {
		return scr.Submit()
	}})

	buttons = append(buttons, components.Button{Title: "[ Back ]", Cmd: func() tea.Cmd {
		return tui.SetBodyPane(tui.WelcomeScreen)
	}})

	scr.inputGroup = components.NewInputGroup(inputs, buttons)

	return scr, err
}

func (s StorageCreateScreen) Init() tea.Cmd {
	return textinput.Blink

}

func (s *StorageCreateScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (s StorageCreateScreen) View() string {
	return screens.RenderContent("Fill in storage details: \n", s.inputGroup.View())
}

func (s *StorageCreateScreen) Submit() tea.Cmd {
	var (
		err  error
		cmds []tea.Cmd
	)

	path := s.inputGroup.Inputs[createPath].Value()
	password := s.inputGroup.Inputs[createPassword].Value()

	// todo: validate inputs

	s.storage, err = s.createStorageUC.Call(path, password, s.encrypter)
	if err != nil {
		cmds = append(cmds, tui.ReportInfo("Error: %v", err))
	} else {
		cmds = append(cmds, tui.ReportInfo("Created new storage: %v", path))
		cmds = append(cmds, tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(s.storage)))
	}

	return tea.Batch(cmds...)
}

type inputOpts struct {
	placeholder string
	charLimit   int
	focus       bool
	value       string
}

func newInput(opts inputOpts) textinput.Model {
	t := textinput.New()
	t.CharLimit = opts.charLimit
	t.Placeholder = opts.placeholder

	if len(opts.value) > 0 {
		t.SetValue(opts.value)
	}

	if opts.focus {
		t.Focus()
		t.PromptStyle = styles.Focused
		t.TextStyle = styles.Focused
	}

	return t
}
