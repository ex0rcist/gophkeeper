package blobedit

import (
	"context"
	"errors"
	"fmt"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/components"
	"gophkeeper/internal/keeper/tui/screens"
	"gophkeeper/internal/keeper/tui/styles"
	"gophkeeper/pkg/models"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	errTitleEmpty    = errors.New("please enter title")
	errMetadataEmpty = errors.New("please enter metadata")
)

const (
	blobTitle = iota
	blobMetadata
)

type BlobEditScreen struct {
	secret  *models.Secret
	storage storage.Storage

	inputGroup components.InputGroup
}

func (s BlobEditScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewBlobEditScreen(msg.Secret, msg.Storage), nil
}

func NewBlobEditScreen(secret *models.Secret, strg storage.Storage) *BlobEditScreen {
	m := BlobEditScreen{
		secret:  secret,
		storage: strg,
	}

	inputs := make([]textinput.Model, 2)
	inputs[blobTitle] = newInput(inputOpts{placeholder: "Title", charLimit: 64})
	inputs[blobMetadata] = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})

	buttons := []components.Button{}
	buttons = append(buttons, components.Button{Title: "[ Pick file ]", Cmd: func() tea.Cmd {

		err := m.validateInputs()
		if err != nil {
			return tui.ReportError(err)
		}

		f := func(args ...any) tea.Cmd {
			str, ok := args[0].(string)
			if !ok {
				return tui.ReportError(fmt.Errorf("error opening file"))
			}

			err := m.Submit(str)
			if err != nil {
				return tui.ReportError(fmt.Errorf("error uploading file: %w", err))
			}

			return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
		}

		return tui.SetBodyPane(tui.FilePickScreen, tui.WithStorage(m.storage), tui.WithCallback(f), tui.WithSecret(secret))
	}})

	buttons = append(buttons, components.Button{Title: "[ Back ]", Cmd: func() tea.Cmd {
		return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
	}})

	if secret.ID > 0 {
		inputs[blobTitle].SetValue(secret.Title)
		inputs[blobMetadata].SetValue(secret.Metadata)
	}

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

func (s BlobEditScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

func (s *BlobEditScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (s BlobEditScreen) View() string {
	return screens.RenderContent("Fill in file details:", s.inputGroup.View())
}

func (s BlobEditScreen) validateInputs() error {
	title := s.inputGroup.Inputs[blobTitle].Value()
	metadata := s.inputGroup.Inputs[blobMetadata].Value()

	if len(title) == 0 {
		return errTitleEmpty
	}

	if len(metadata) == 0 {
		return errMetadataEmpty
	}

	return nil
}

func (s *BlobEditScreen) Submit(path string) error {
	var (
		err error
	)

	err = s.validateInputs()
	if err != nil {
		return err
	}

	bts, err := readFileToBytes(path)
	if err != nil {
		return err
	}

	s.secret.Title = s.inputGroup.Inputs[blobTitle].Value()
	s.secret.Metadata = s.inputGroup.Inputs[blobMetadata].Value()
	s.secret.Blob = &models.Blob{
		FileName:  path,
		FileBytes: bts,
	}
	s.secret.UpdatedAt = time.Now()

	// Save secret
	if s.secret.ID == 0 {
		s.secret.CreatedAt = time.Now()
		err = s.storage.Create(context.Background(), s.secret)
	} else {
		err = s.storage.Update(context.Background(), s.secret)
	}

	return err
}

func readFileToBytes(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
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
