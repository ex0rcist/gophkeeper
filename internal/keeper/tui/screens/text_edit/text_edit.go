package textedit

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
	errContentEmpty  = errors.New("Please enter content")
	errTitleEmpty    = errors.New("Please enter title")
)

const (
	textTitle = iota
	textMetadata
	textContent
)

type TextEditScreen struct {
	secret  *models.Secret
	storage storage.Storage

	inputGroup components.InputGroup
}

func (s TextEditScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewTextEditScreen(msg.Secret, msg.Storage), nil
}

func NewTextEditScreen(secret *models.Secret, strg storage.Storage) *TextEditScreen {
	m := TextEditScreen{
		secret:  secret,
		storage: strg,
	}

	inputs := make([]textinput.Model, 3)
	inputs[textTitle] = newInput(inputOpts{placeholder: "Title", charLimit: 64})
	inputs[textMetadata] = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})
	inputs[textContent] = newInput(inputOpts{placeholder: "Content", charLimit: 164})

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
		inputs[textTitle].SetValue(secret.Title)
		inputs[textMetadata].SetValue(secret.Metadata)
		inputs[textContent].SetValue(secret.Text.Content)

	}

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

func (s TextEditScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

func (s *TextEditScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Handle input group. TODO: fix blink
	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (s *TextEditScreen) Submit() error {
	var (
		err error
	)

	title := s.inputGroup.Inputs[textTitle].Value()
	metadata := s.inputGroup.Inputs[textMetadata].Value()
	content := s.inputGroup.Inputs[textContent].Value()

	// Validate inputs
	if len(metadata) == 0 {
		return errMetadataEmpty
	}

	if len(title) == 0 {
		return errTitleEmpty
	}

	if len(content) == 0 {
		return errContentEmpty
	}

	s.secret.Title = title
	s.secret.Metadata = metadata
	s.secret.Text = &models.Text{Content: content}
	s.secret.UpdatedAt = time.Now()

	// Save text
	if s.secret.ID == 0 {
		s.secret.CreatedAt = time.Now()
		err = s.storage.Create(context.Background(), *s.secret)
	} else {
		err = s.storage.Update(context.Background(), *s.secret)
	}

	return err
}

func (s TextEditScreen) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Fill in text details: \n"))
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
