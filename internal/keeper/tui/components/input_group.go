// Provides different building blocks for TUI
package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textinput"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()
)

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
			input.PromptStyle = focusedStyle
			input.TextStyle = focusedStyle
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
					m.Inputs[i].PromptStyle = focusedStyle
					m.Inputs[i].TextStyle = focusedStyle
					continue
				}

				// Remove focused state
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = noStyle
				m.Inputs[i].TextStyle = noStyle
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
		b     strings.Builder
		style lipgloss.Style = blurredStyle
	)

	for i, input := range m.Inputs {
		// TODO
		b.WriteString(fmt.Sprintf("%s: %s", input.Placeholder, input.View()))

		if i < m.totalPos-1 {
			b.WriteRune('\n')
		}
	}

	for i, but := range m.Buttons {
		title := but.Title

		if m.FocusIndex == len(m.Inputs)+i {
			style = focusedStyle
		} else {
			style = blurredStyle
		}

		if i < m.totalPos-1 {
			b.WriteRune('\n')
		}

		fmt.Fprintf(&b, "%s", style.Render(title))
	}

	return b.String()
}
