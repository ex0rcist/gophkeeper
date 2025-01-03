package blobedit

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"gophkeeper/internal/keeper/storage"
// 	"gophkeeper/internal/keeper/tui"
// 	"gophkeeper/pkg/models"

// 	"os"
// 	"strings"

// 	"github.com/charmbracelet/bubbles/filepicker"
// 	"github.com/charmbracelet/bubbles/textinput"
// 	"github.com/charmbracelet/lipgloss"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// type status int

// const (
// 	// Statuses
// 	fileStart status = iota
// 	fileStartDownload
// 	filePicking
// 	fileUpload
// 	fileDownload
// 	fileComplete
// 	fileError
// )

// var (
// 	errMetadataEmpty = errors.New("Please enter metadata")
// )

// var (
// 	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
// 	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
// 	cursorStyle         = focusedStyle
// 	noStyle             = lipgloss.NewStyle()
// 	helpStyle           = blurredStyle
// 	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

// 	focusedButton = focusedStyle.Render("[ Submit ]")
// 	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
// )

// type fileCompleteMsg struct{}

// type errMsg struct {
// 	msg string
// }

// // Model for uploading or downloading file from server
// type FileEditScreen struct {
// 	secret  *models.Secret
// 	storage storage.Storage

// 	metadataInput textinput.Model
// 	filepicker    filepicker.Model
// 	status        status
// 	isDownload    bool
// 	selectedFile  string
// }

// type Button struct {
// 	Title string
// 	Cmd   func() tea.Cmd
// }

// func (s FileEditScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
// 	return NewFileEditScreen(msg.Secret, msg.Storage), nil
// }

// func NewFileEditScreen(secret *models.Secret, strg storage.Storage) *FileModel {
// 	var err error

// 	m := FileEditScreen{
// 		metadataInput: newInput(inputOpts{placeholder: "Metadata", charLimit: 64}),
// 		secret:        secret,
// 		storage:       strg,
// 		status:        fileStart,
// 	}

// 	fp := filepicker.New()
// 	fp.AutoHeight = false

// 	fp.CurrentDirectory, err = os.UserHomeDir()
// 	if err != nil {
// 		// todo report error
// 	}

// 	// Set download mode if ID is passed
// 	if secret.ID > 0 {
// 		m.isDownload = true
// 		m.status = fileStartDownload
// 	}

// 	m.filepicker = fp

// 	return &m
// }

// func (m FileEditScreen) Init() tea.Cmd {
// 	return textinput.Blink
// }

// func (m *FileEditScreen) Update(msg tea.Msg) tea.Cmd {
// 	var (
// 		cmds []tea.Cmd
// 		cmd  tea.Cmd
// 	)

// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {

// 		case "p":
// 			if m.status == fileStart {
// 				// Unfocus metadata
// 				m.metadataInput.Blur()
// 				// m.metadataInput.PromptStyle = style.BlurredStyle
// 				// m.metadataInput.TextStyle = style.BlurredStyle

// 				m.status = filePicking
// 				return m.filepicker.Init()
// 			}
// 		case "b":
// 			if m.status == filePicking {
// 				// Focus metadata input
// 				// m.metadataInput.PromptStyle = style.FocusedStyle
// 				// m.metadataInput.TextStyle = style.FocusedStyle

// 				m.status = fileStart
// 				return m.metadataInput.Focus()
// 			}
// 		case "d":
// 			if m.status == fileStartDownload {
// 				m.status = fileDownload

// 				return m.downloadStart()
// 			}
// 		}
// 	case fileCompleteMsg:
// 		m.status = fileComplete
// 	case errMsg:
// 		// m.errorMsg = msg.msg // todo report
// 		m.status = fileError
// 	}

// 	// Handle metadata input
// 	if m.status == fileStart {
// 		m.metadataInput, cmd = m.metadataInput.Update(msg)
// 		cmds = append(cmds, cmd)
// 	}

// 	// Set filepicker size
// 	m.filepicker.Height = 30 // m.state.windowHeight - 10 // TODO

// 	// Update file picker if it's in focus
// 	if m.status == filePicking {
// 		m.filepicker, cmd = m.filepicker.Update(msg)
// 		cmds = append(cmds, cmd)

// 		// Upload file if user picked a file
// 		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
// 			m.selectedFile = path

// 			return m.uploadStart()
// 		}
// 	}

// 	return tea.Batch(cmds...)
// }

// func (m *FileEditScreen) uploadStart() tea.Cmd {
// 	m.status = fileUpload

// 	return func() tea.Msg {
// 		// validate inputs
// 		metadata := m.metadataInput.Value()
// 		if len(metadata) == 0 {
// 			return errMsg{msg: errMetadataEmpty.Error()}
// 		}

// 		m.secret.Metadata = metadata
// 		m.secret.Blob = &models.Blob{FileName: "test.wtf"}

// 		// upload file
// 		err := m.storage.Create(context.Background(), *m.secret)
// 		if err != nil {
// 			return errMsg{msg: err.Error()}
// 		}

// 		return fileCompleteMsg{}
// 	}
// }

// func (m *FileEditScreen) downloadStart() tea.Cmd {
// 	return func() tea.Msg {
// 		// load secret
// 		// secret, err := m.state.client.LoadSecret(context.Background(), m.secretID)
// 		// if err != nil {
// 		// 	return errMsg{msg: err.Error()}
// 		// }

// 		// // download file
// 		// err = m.state.client.DownloadFile(context.Background(), m.secretID, secret.Blob.FileName)
// 		// if err != nil {
// 		// 	return errMsg{msg: err.Error()}
// 		// }

// 		return fileCompleteMsg{}
// 	}
// }

// func (m FileEditScreen) View() string {
// 	var b strings.Builder

// 	if m.isDownload {
// 		// Download mode
// 		b.WriteString("Press d to start file download or esc to go back\n")
// 	} else {
// 		// Upload mode
// 		b.WriteString("Pick a file and enter metadata\n")

// 		b.WriteString(m.metadataInput.View())
// 		b.WriteString("\n\n")
// 	}

// 	switch m.status {
// 	case fileStart:
// 		b.WriteString("Press p to pick file")
// 	case filePicking:
// 		b.WriteString("Press b to edit metadata\n")
// 		b.WriteString(m.filepicker.View())
// 	case fileUpload:
// 		b.WriteString("File upload in progress, please wait...")
// 	case fileDownload:
// 		b.WriteString("File download in progress, please wait...")
// 	case fileComplete:
// 		b.WriteString("File transfer is done, press esc to go back")
// 	case fileError:
// 		b.WriteString("Error occured during file transfer")
// 	}

// 	body := style.RenderBox(b.String())

// 	if len(m.errorMsg) > 0 {
// 		errorBox := style.ErrorStyle.Render(m.errorMsg)

// 		body = lipgloss.JoinVertical(lipgloss.Top, body, errorBox)
// 	}

// 	return body
// }

// type inputOpts struct {
// 	placeholder string
// 	charLimit   int
// 	focus       bool
// }

// func newInput(opts inputOpts) textinput.Model {
// 	t := textinput.New()
// 	t.CharLimit = opts.charLimit
// 	t.Placeholder = opts.placeholder

// 	if opts.focus {
// 		t.Focus()
// 		t.PromptStyle = focusedStyle
// 		t.TextStyle = focusedStyle
// 	}

// 	return t
// }

//// -========== -========== -========== -========== -========== -========== -========== -========== -========== -========== -========== -==========

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"gophkeeper/internal/keeper/storage"
// 	"gophkeeper/internal/keeper/tui"
// 	"gophkeeper/pkg/models"
// 	"strings"
// 	"time"

// 	"github.com/charmbracelet/bubbles/textarea"
// 	"github.com/charmbracelet/bubbles/textinput"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// var (
// 	errMetadataEmpty = errors.New("Please enter metadata")
// 	errTextEmpty     = errors.New("Please enter login")
// 	errTitleEmpty    = errors.New("Please enter title")
// )

// const (
// 	textTitle = iota
// 	textMetadata
// 	textContent
// )

// var (
// 	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
// 	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
// 	cursorStyle         = focusedStyle
// 	noStyle             = lipgloss.NewStyle()
// 	helpStyle           = blurredStyle
// 	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

// 	focusedButton = focusedStyle.Render("[ Submit ]")
// 	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
// )

// type TextEditScreen struct {
// 	secret  *models.Secret
// 	storage storage.Storage

// 	textTitleInput    textinput.Model
// 	textMetadataInput textinput.Model
// 	textContentInput  textarea.Model

// 	submitButton Button
// 	backButton   Button

// 	focusIndex int
// 	totalPos   int
// }

// type Button struct {
// 	Title string
// 	Cmd   func() tea.Cmd
// }

// func (s TextEditScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
// 	return NewTextEditScreen(msg.Secret, msg.Storage), nil
// }

// func NewTextEditScreen(secret *models.Secret, strg storage.Storage) *TextEditScreen {
// 	m := TextEditScreen{
// 		secret:   secret,
// 		storage:  strg,
// 		totalPos: 4,
// 	}

// 	m.textTitleInput = newInput(inputOpts{placeholder: "Title", charLimit: 64})
// 	m.textMetadataInput = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})
// 	m.textContentInput = newTextareaInput(inputOpts{placeholder: "Text", charLimit: 64})

// 	m.submitButton = Button{Title: "[ Submit ]", Cmd: func() tea.Cmd {
// 		err := m.Submit()
// 		if err != nil {
// 			return tui.ReportError(err)
// 		} else {
// 			// todo: invalidate or update cache ?
// 			return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
// 		}
// 	}}

// 	m.backButton = Button{Title: "[ Back ]", Cmd: func() tea.Cmd {
// 		return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
// 	}}

// 	if secret.ID > 0 {
// 		m.textTitleInput.SetValue(secret.Title)
// 		m.textMetadataInput.SetValue(secret.Metadata)
// 		m.textContentInput.SetValue(secret.Text.Content)
// 	}

// 	return &m
// }

// func (s TextEditScreen) Init() tea.Cmd {
// 	return nil // s.inputGroup.Init() ???
// }

// func (s *TextEditScreen) Update(msg tea.Msg) tea.Cmd {
// 	var (
// 		cmd  tea.Cmd
// 		cmds []tea.Cmd
// 	)

// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "enter", "up", "down":
// 			str := msg.String()

// 			if str == "ctrl+enter" && s.focusIndex == 2 {
// 				s.textContentInput.
// 			}

// 			// Did the user press enter while the button was focused?
// 			if str == "enter" {

// 				if s.focusIndex == 3 {
// 					return s.submitButton.Cmd()
// 				}

// 				if s.focusIndex == 4 {
// 					return s.backButton.Cmd()
// 				}
// 			}

// 			// Cycle indexes
// 			if str == "up" {
// 				s.focusIndex--
// 			} else {
// 				s.focusIndex++
// 			}

// 			if s.focusIndex > s.totalPos {
// 				s.focusIndex = 0
// 			} else if s.focusIndex < 0 {
// 				s.focusIndex = s.totalPos
// 			}

// 			if s.focusIndex == 0 {
// 				cmds = append(cmds, s.textTitleInput.Focus())
// 				s.textTitleInput.PromptStyle = focusedStyle
// 				s.textTitleInput.TextStyle = focusedStyle
// 			} else {
// 				s.textTitleInput.Blur()
// 				s.textTitleInput.PromptStyle = noStyle
// 				s.textTitleInput.TextStyle = noStyle
// 			}

// 			if s.focusIndex == 1 {
// 				cmds = append(cmds, s.textMetadataInput.Focus())
// 				s.textMetadataInput.PromptStyle = focusedStyle
// 				s.textMetadataInput.TextStyle = focusedStyle
// 			} else {
// 				s.textMetadataInput.Blur()
// 				s.textMetadataInput.PromptStyle = noStyle
// 				s.textMetadataInput.TextStyle = noStyle
// 			}

// 			if s.focusIndex == 2 {
// 				cmds = append(cmds, s.textContentInput.Focus())
// 				// s.textContentInput.PromptStyle = focusedStyle
// 				// s.textContentInput.TextStyle = focusedStyle
// 			} else {
// 				s.textContentInput.Blur()
// 				// s.textMetadataInput.PromptStyle = noStyle
// 				// s.textMetadataInput.TextStyle = noStyle
// 			}

// 		}
// 	}

// 	mm, cmd := s.textTitleInput.Update(msg)
// 	s.textTitleInput = mm
// 	cmds = append(cmds, cmd)

// 	mm, cmd = s.textMetadataInput.Update(msg)
// 	s.textMetadataInput = mm
// 	cmds = append(cmds, cmd)

// 	mmm, cmd := s.textContentInput.Update(msg)
// 	s.textContentInput = mmm
// 	cmds = append(cmds, cmd)

// 	// for i := range m.Inputs {
// 	// 	m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
// 	// }

// 	return tea.Batch(cmds...)
// }

// func (s *TextEditScreen) Submit() error {
// 	var (
// 		err error
// 	)

// 	title := s.textTitleInput.Value()
// 	metadata := s.textMetadataInput.Value()
// 	text := s.textContentInput.Value()

// 	if len(metadata) == 0 {
// 		return errMetadataEmpty
// 	}

// 	if len(title) == 0 {
// 		return errTitleEmpty
// 	}

// 	if len(text) == 0 {
// 		return errTextEmpty
// 	}

// 	s.secret.Title = title
// 	s.secret.Metadata = metadata
// 	s.secret.Text = &models.Text{Content: text}
// 	s.secret.UpdatedAt = time.Now()

// 	// Save credential
// 	if s.secret.ID == 0 {
// 		s.secret.CreatedAt = time.Now()
// 		err = s.storage.Create(context.Background(), *s.secret)
// 	} else {
// 		err = s.storage.Update(context.Background(), *s.secret)
// 	}

// 	return err
// }

// func (s TextEditScreen) View() string {
// 	var (
// 		b strings.Builder
// 		// style lipgloss.Style = blurredStyle
// 	)

// 	b.WriteString(fmt.Sprintf("Fill in text details: %d/%d \n", s.focusIndex, s.totalPos))

// 	b.WriteString(fmt.Sprintf("%s\n", s.textTitleInput.View()))
// 	b.WriteString(fmt.Sprintf("%s\n", s.textMetadataInput.View()))
// 	b.WriteString(fmt.Sprintf("%s\n", s.textContentInput.View()))

// 	if s.focusIndex == 3 {
// 		b.WriteString(focusedStyle.Render(fmt.Sprintf("%s\n", s.submitButton.Title)))
// 	} else {
// 		b.WriteString(blurredStyle.Render(fmt.Sprintf("%s\n", s.submitButton.Title)))
// 	}

// 	if s.focusIndex == 4 {
// 		b.WriteString(focusedStyle.Render(fmt.Sprintf("%s\n", s.backButton.Title)))
// 	} else {
// 		b.WriteString(blurredStyle.Render(fmt.Sprintf("%s\n", s.backButton.Title)))
// 	}

// 	return b.String()
// }

// type inputOpts struct {
// 	placeholder string
// 	charLimit   int
// 	focus       bool
// }

// func newInput(opts inputOpts) textinput.Model {
// 	t := textinput.New()
// 	t.CharLimit = opts.charLimit
// 	t.Placeholder = opts.placeholder

// 	if opts.focus {
// 		t.Focus()
// 		t.PromptStyle = focusedStyle
// 		t.TextStyle = focusedStyle
// 	}

// 	return t
// }

// func newTextareaInput(opts inputOpts) textarea.Model {
// 	t := textarea.New()
// 	t.Placeholder = "Tell me all your secrets!"
// 	if opts.focus {
// 		t.Focus()
// 		// ?
// 	}

// 	return t
// }
