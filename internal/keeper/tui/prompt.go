package tui

import (
	"fmt"
	"gophkeeper/internal/keeper/tui/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PromptMsg enables the prompt widget.
type PromptMsg struct {
	Prompt       string       // Prompt to display to the user
	Placeholder  string       // Set placeholder text in prompt
	InitialValue string       // Set initial value for the user to edit
	Action       PromptAction // Action to carry out when key is pressed
	Key          key.Binding  // Key that when pressed triggers the action and closes the prompt
	Cancel       key.Binding  // Cancel is a key that when pressed skips the action and closes the prompt
	AnyCancel    bool         // If any key can cancel the prompt
}

type PromptAction func(text string) tea.Cmd

func StringPrompt(prompt string, action PromptAction) tea.Cmd {
	return CmdHandler(PromptMsg{
		Prompt: fmt.Sprintf("%s: ", prompt),
		Action: action,
		Key: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		AnyCancel: false,
	})
}

// Yes/No question. If yes is given then the action is invoked.
func YesNoPrompt(prompt string, action tea.Cmd) tea.Cmd {
	return CmdHandler(PromptMsg{
		Prompt: fmt.Sprintf("%s (y/N): ", prompt),
		Action: func(_ string) tea.Cmd {
			return action
		},
		Key: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "confirm"),
		),
		AnyCancel: true,
	})
}

func NewPrompt(msg PromptMsg) (*Prompt, tea.Cmd) {
	model := textinput.New()
	model.Prompt = msg.Prompt
	model.SetValue(msg.InitialValue)
	model.Placeholder = msg.Placeholder
	model.PlaceholderStyle = styles.Regular.Faint(true)
	blink := model.Focus()

	prompt := Prompt{
		model:     model,
		action:    msg.Action,
		trigger:   msg.Key,
		cancel:    msg.Cancel,
		anyCancel: msg.AnyCancel,
	}
	return &prompt, blink
}

// Prompt is a widget that prompts the user for input and triggers an action.
type Prompt struct {
	model     textinput.Model
	action    PromptAction
	trigger   key.Binding
	cancel    key.Binding
	anyCancel bool
}

// HandleKey handles the user key press, and returns a command to be run, and
// whether the prompt should be closed.
func (p *Prompt) HandleKey(msg tea.KeyMsg) (closePrompt bool, cmd tea.Cmd) {
	switch {
	case key.Matches(msg, p.trigger):
		cmd = p.action(p.model.Value())
		closePrompt = true
	case key.Matches(msg, p.cancel), p.anyCancel:
		cmd = ReportInfo("canceled operation")
		closePrompt = true
	default:
		p.model, cmd = p.model.Update(msg)
	}
	return
}

// HandleBlink handles the bubbletea blink message.
func (p *Prompt) HandleBlink(msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Ignore key presses, they're handled by HandleKey above.
	default:
		// The blink message type is unexported so we just send unknown types to the model.
		p.model, cmd = p.model.Update(msg)
	}
	return
}

func (p *Prompt) View(width int) string {
	paddedBorder := styles.ThickBorder.BorderForeground(styles.Red).Padding(0, 1)
	paddedBorderWidth := paddedBorder.GetHorizontalBorderSize() + paddedBorder.GetHorizontalPadding()
	// Set available width for user entered value before it horizontally
	// scrolls.
	p.model.Width = max(0, width-lipgloss.Width(p.model.Prompt)-paddedBorderWidth)
	// Render a prompt, surrounded by a padded red border, spanning the width of the
	// terminal, accounting for width of border. Inline and MaxWidth ensures the
	// prompt remains on a single line.
	content := styles.Regular.Inline(true).MaxWidth(width - paddedBorderWidth).Render(p.model.View())
	return paddedBorder.Width(width - paddedBorder.GetHorizontalBorderSize()).Render(content)
}

func (p *Prompt) HelpBindings() []key.Binding {
	bindings := []key.Binding{
		p.trigger,
	}
	if p.anyCancel {
		bindings = append(bindings, key.NewBinding(key.WithHelp("n", "cancel")))
	} else {
		bindings = append(bindings, p.cancel)
	}
	return bindings
}
