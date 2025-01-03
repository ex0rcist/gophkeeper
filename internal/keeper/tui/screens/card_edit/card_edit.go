package cardedit

import (
	"context"
	"errors"
	"fmt"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/components"
	"gophkeeper/pkg/models"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	errFieldEmpty = func(label string) error { return errors.New(fmt.Sprintf("Please enter %s", label)) }
)

const (
	cardTitle = iota
	cardMetadata
	cardNumber
	cardExpYear
	cardExpMonth
	cardCVV
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type CardEditScreen struct {
	secret  *models.Secret
	storage storage.Storage

	inputGroup components.InputGroup
}

func (s CardEditScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewCardEditScreen(msg.Secret, msg.Storage), nil
}

func NewCardEditScreen(secret *models.Secret, strg storage.Storage) *CardEditScreen {
	m := CardEditScreen{
		secret:  secret,
		storage: strg,
	}

	inputs := make([]textinput.Model, 6)
	inputs[cardTitle] = newInput(inputOpts{placeholder: "Title", charLimit: 64})
	inputs[cardMetadata] = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})
	inputs[cardNumber] = newInput(inputOpts{placeholder: "Card number", charLimit: 64})
	inputs[cardExpYear] = newInput(inputOpts{placeholder: "Exp Year", charLimit: 2})
	inputs[cardExpMonth] = newInput(inputOpts{placeholder: "Exp Month", charLimit: 2})
	inputs[cardCVV] = newInput(inputOpts{placeholder: "CVV", charLimit: 6})

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
		inputs[cardTitle].SetValue(secret.Title)
		inputs[cardMetadata].SetValue(secret.Metadata)
		inputs[cardNumber].SetValue(secret.Card.Number)
		inputs[cardExpMonth].SetValue(strconv.FormatUint(uint64(secret.Card.ExpMonth), 10))
		inputs[cardExpYear].SetValue(strconv.FormatUint(uint64(secret.Card.ExpYear), 10))
		inputs[cardCVV].SetValue(strconv.FormatUint(uint64(secret.Card.CVV), 10))
	}

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

func (s CardEditScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

func (s *CardEditScreen) Update(msg tea.Msg) tea.Cmd {
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

func (s *CardEditScreen) Submit() error {
	var (
		err error
	)

	title := s.inputGroup.Inputs[cardTitle].Value()
	metadata := s.inputGroup.Inputs[cardMetadata].Value()
	cardNumber := s.inputGroup.Inputs[cardNumber].Value()
	cardExpMonth := s.inputGroup.Inputs[cardExpMonth].Value()
	cardExpYear := s.inputGroup.Inputs[cardExpYear].Value()
	cardCVV := s.inputGroup.Inputs[cardCVV].Value()

	// Validate inputs
	if len(metadata) == 0 {
		return errFieldEmpty("metadata")
	}

	if len(title) == 0 {
		return errFieldEmpty("title")
	}

	if len(cardNumber) == 0 {
		return errFieldEmpty("card number")
	}

	if len(cardNumber) == 0 {
		return errFieldEmpty("card number")
	}

	if len(cardExpYear) == 0 {
		return errFieldEmpty("exp year")
	}

	if len(cardExpMonth) == 0 {
		return errFieldEmpty("exp month")
	}

	if len(cardCVV) == 0 {
		return errFieldEmpty("CVV")
	}

	s.secret.Title = title
	s.secret.Metadata = metadata
	card := &models.Card{Number: cardNumber}

	card.ExpYear = strToUint32(cardExpYear)
	card.ExpMonth = strToUint32(cardExpMonth)
	card.CVV = strToUint32(cardCVV)

	s.secret.Card = card
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

func (s CardEditScreen) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Fill in credential details: \n"))
	b.WriteString(s.inputGroup.View())

	// 	cardForm := fmt.Sprintf(
	// 		`
	// %s
	// %s

	// %s %s %s
	// %s %s %s
	// `,
	// 		style.FocusedStyle.Width(30).Render("Card Number"),
	// 		m.inputGroup.Inputs[cardNumber].View(),
	// 		style.FocusedStyle.Width(8).Render("Exp MM"),
	// 		style.FocusedStyle.Width(8).Render("Exp YY"),
	// 		style.FocusedStyle.Width(6).Render("CVV"),
	// 		m.inputGroup.Inputs[cardExpMonth].View(),
	// 		m.inputGroup.Inputs[cardExpYear].View(),
	// 		m.inputGroup.Inputs[cardCvv].View(),
	// 	)

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
		t.PromptStyle = focusedStyle
		t.TextStyle = focusedStyle
	}

	return t
}

func strToUint32(str string) uint32 {
	i64, _ := strconv.ParseUint(str, 10, 32)
	return uint32(i64)
}
