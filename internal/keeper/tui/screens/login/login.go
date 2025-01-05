package login

import (
	"context"
	"errors"
	"gophkeeper/internal/keeper/api"
	"gophkeeper/internal/keeper/crypto"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/components"
	"gophkeeper/internal/keeper/tui/screens"
	"gophkeeper/internal/keeper/tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	errLoginEmpty    = errors.New("Please enter login")
	errPasswordEmpty = errors.New("Please enter password")
)

const (
	posLogin = iota
	posPassword
)

type Mode int

const (
	modeLogin Mode = iota
	modeRegister
)

type LoginScreen struct {
	client    api.IApiClient
	encrypter crypto.Encrypter

	inputGroup components.InputGroup
}

func (m LoginScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewLoginScreen(msg.Client), nil
}

func NewLoginScreen(client api.IApiClient) *LoginScreen {
	m := LoginScreen{
		client: client,
	}

	inputs := make([]textinput.Model, 2)
	inputs[posLogin] = newInput(inputOpts{placeholder: "Login", charLimit: 64})
	inputs[posPassword] = newInput(inputOpts{placeholder: "Password", charLimit: 64})

	buttons := []components.Button{}
	buttons = append(buttons, components.Button{Title: "[ Login ]", Cmd: func() tea.Cmd {
		return m.Submit(modeLogin)
	}})

	buttons = append(buttons, components.Button{Title: "[ Register ]", Cmd: func() tea.Cmd {
		return m.Submit(modeRegister)
	}})

	buttons = append(buttons, components.Button{Title: "[ Back ]", Cmd: func() tea.Cmd {
		return tui.GoToStart()
	}})

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

func (s LoginScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

func (s *LoginScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (s *LoginScreen) Submit(mode Mode) tea.Cmd {
	var (
		token string
		err   error
		cmds  []tea.Cmd
	)

	login := s.inputGroup.Inputs[posLogin].Value()
	password := s.inputGroup.Inputs[posPassword].Value()

	// Validate inputs
	if len(login) == 0 {
		return tui.ReportError(errLoginEmpty)
	}
	if len(password) == 0 {
		return tui.ReportError(errPasswordEmpty)
	}

	switch mode {
	case modeLogin:
		token, err = s.client.Login(context.Background(), login, password)
	case modeRegister:
		token, err = s.client.Register(context.Background(), login, password)
	}

	if err != nil {
		cmds = append(cmds, tui.ReportError(err))
	} else {
		s.client.SetToken(token)
		s.client.SetPassword(password)

		// create storage instance
		storage, err := storage.NewRemoteStorage(s.client, s.encrypter)
		if err != nil {
			cmds = append(cmds, tui.ReportError(err))
		} else {
			cmds = append(cmds, tui.ReportInfo("success!"))
			cmds = append(cmds, tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(storage)))
		}
	}

	return tea.Batch(cmds...)
}

func (s LoginScreen) View() string {
	return screens.RenderContent("Fill in credentials:", s.inputGroup.View())
}

type inputOpts struct {
	placeholder string
	charLimit   int
	focus       bool
}

func newInput(opts inputOpts) textinput.Model {
	t := textinput.New()
	t.CharLimit = opts.charLimit
	t.Placeholder = opts.placeholder

	if opts.focus {
		t.Focus()
		t.PromptStyle = styles.Focused
		t.TextStyle = styles.Focused
	}

	return t
}
