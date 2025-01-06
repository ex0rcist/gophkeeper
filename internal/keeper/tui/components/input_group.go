// Provides different building blocks for TUI
package components

import (
	"fmt"
	"gophkeeper/internal/keeper/tui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textinput"
)

var ()

type InputGroup struct {
	Inputs     []textinput.Model
	Buttons    []Button
	FocusIndex int

	totalPos int // total positions for cursor
}

type Button struct {
	Title string
	Cmd   func() tea.Cmd
}

func NewInputGroup(inputs []textinput.Model, buttons []Button) InputGroup {
	// Set styles
	for i, input := range inputs {
		if i == 0 {
			input.Focus()
			input.PromptStyle = styles.Focused
			input.TextStyle = styles.Focused
		}

		inputs[i] = input
	}

	return InputGroup{
		Inputs:   inputs,
		Buttons:  buttons,
		totalPos: len(inputs) + len(buttons) - 1,
	}
}

func (m InputGroup) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the button was focused?
			if s == "enter" {
				butIdx := m.FocusIndex - len(m.Inputs)
				if butIdx >= 0 {
					return m, m.Buttons[butIdx].Cmd()
				}
			}

			// Cycle indexes
			if s == "up" {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}

			if m.FocusIndex > m.totalPos {
				m.FocusIndex = 0
			} else if m.FocusIndex < 0 {
				m.FocusIndex = m.totalPos
			}

			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds = append(cmds, m.Inputs[i].Focus())
					m.Inputs[i].PromptStyle = styles.Focused
					m.Inputs[i].TextStyle = styles.Focused
					continue
				}

				// Remove focused state
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = styles.Regular
				m.Inputs[i].TextStyle = styles.Regular
			}

		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *InputGroup) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m InputGroup) View() string {
	var (
		b       strings.Builder
		style   lipgloss.Style
		padding int // num spaces to pad
	)

	// Calc max label width
	maxLabelLength := 0
	for _, input := range m.Inputs {
		if len(input.Placeholder) > maxLabelLength {
			maxLabelLength = len(input.Placeholder)
		}
	}

	// Draw inputs
	for _, input := range m.Inputs {
		label := input.Placeholder
		padding = maxLabelLength - len(label) // Align right

		b.WriteString(fmt.Sprintf("%s: %s\n",
			strings.Repeat(" ", padding)+label, // Add padding
			input.View(),
		))
	}

	b.WriteRune('\n')

	// Draw buttons
	buttonPadding := maxLabelLength + 2 // 2 for `: `
	for i, but := range m.Buttons {
		title := but.Title

		if m.FocusIndex == len(m.Inputs)+i {
			style = styles.Focused
		} else {
			style = styles.Blurred
		}

		b.WriteString(fmt.Sprintf("%s%s\n",
			strings.Repeat(" ", buttonPadding),
			style.Render(title),
		))
	}

	return b.String()
}
