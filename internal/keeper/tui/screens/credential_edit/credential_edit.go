package credentialedit

import (
	"context"
	"errors"
	"fmt"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/components"
	"gophkeeper/internal/keeper/tui/styles"
	"gophkeeper/pkg/models"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	errMetadataEmpty = errors.New("Please enter metadata")
	errLoginEmpty    = errors.New("Please enter login")
	errPasswordEmpty = errors.New("Please enter password")
	errTitleEmpty    = errors.New("Please enter title")
)

const (
	credTitle = iota
	credMetadata
	credLogin
	credPassword
)

type CredentialEditScreen struct {
	secret  *models.Secret
	storage storage.Storage

	inputGroup components.InputGroup
}

func (s CredentialEditScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewCredentialEditScreen(msg.Secret, msg.Storage), nil
}

func NewCredentialEditScreen(secret *models.Secret, strg storage.Storage) *CredentialEditScreen {
	m := CredentialEditScreen{
		secret:  secret,
		storage: strg,
	}

	inputs := make([]textinput.Model, 4)
	inputs[credTitle] = newInput(inputOpts{placeholder: "Title", charLimit: 64})
	inputs[credMetadata] = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})
	inputs[credLogin] = newInput(inputOpts{placeholder: "Login", charLimit: 64})
	inputs[credPassword] = newInput(inputOpts{placeholder: "Password", charLimit: 64})

	buttons := []components.Button{}
	buttons = append(buttons, components.Button{Title: "[ Submit ]", Cmd: func() tea.Cmd {
		err := m.Submit()
		log.Println(err)
		if err != nil {
			return tui.ReportError(err)
		} else {
			// todo: invalidate or update cache ?
			return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
		}
	}})

	buttons = append(buttons, components.Button{Title: "[ Back ]", Cmd: func() tea.Cmd {
		return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
	}})

	if secret.ID > 0 {
		inputs[credTitle].SetValue(secret.Title)
		inputs[credMetadata].SetValue(secret.Metadata)
		inputs[credLogin].SetValue(secret.Creds.Login)
		inputs[credPassword].SetValue(secret.Creds.Password)
	}

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

func (s CredentialEditScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

func (s *CredentialEditScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":

		}
	}

	// Handle input group. TODO: fix blink
	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (s *CredentialEditScreen) Submit() error {
	var (
		err error
	)

	title := s.inputGroup.Inputs[credTitle].Value()
	metadata := s.inputGroup.Inputs[credMetadata].Value()
	login := s.inputGroup.Inputs[credLogin].Value()
	password := s.inputGroup.Inputs[credPassword].Value()

	// Validate inputs
	if len(metadata) == 0 {
		return errMetadataEmpty
	}

	if len(title) == 0 {
		return errTitleEmpty
	}

	if len(login) == 0 {
		return errLoginEmpty
	}

	if len(password) == 0 {
		return errPasswordEmpty
	}

	s.secret.Title = title
	s.secret.Metadata = metadata
	s.secret.Creds = &models.Credentials{Login: login, Password: password}
	s.secret.UpdatedAt = time.Now()

	// Save credential
	if s.secret.ID == 0 {
		s.secret.CreatedAt = time.Now()
		err = s.storage.Create(context.Background(), *s.secret)
	} else {
		err = s.storage.Update(context.Background(), *s.secret)
	}

	return err
}

func (s CredentialEditScreen) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Fill in credential details: \n"))
	b.WriteString(s.inputGroup.View())

	return b.String()
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
