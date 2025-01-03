package top

import (
	"fmt"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/keys"
	"gophkeeper/internal/keeper/tui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	normalMode mode = iota // default
	promptMode             // confirm prompt is visible and taking input
)

type model struct {
	*tui.PaneManager

	makers map[tui.Screen]tui.ScreenMaker
	// modules    *module.Service
	width    int
	height   int
	mode     mode
	showHelp bool
	prompt   *tui.Prompt

	err  error
	info string
}

func newModel() (model, error) {

	makers := makeMakers()

	m := model{
		PaneManager: tui.NewPaneManager(makers),
		makers:      makers,
	}
	return m, nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.PaneManager.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// if m.dump != nil {
	// 	spew.Fdump(m.dump, msg)
	// }

	switch msg := msg.(type) {
	case tui.PromptMsg:
		// Enable prompt widget
		m.mode = promptMode
		var blink tea.Cmd
		m.prompt, blink = tui.NewPrompt(msg)
		// Send out message to panes to resize themselves to make room for the prompt above it.
		_ = m.PaneManager.Update(tea.WindowSizeMsg{
			Height: m.viewHeight(),
			Width:  m.viewWidth(),
		})
		return m, tea.Batch(cmd, blink)
	case tea.KeyMsg:
		// Pressing any key makes any info/error message in the footer disappear
		m.info = ""
		m.err = nil

		switch m.mode {
		case promptMode:
			closePrompt, cmd := m.prompt.HandleKey(msg)
			if closePrompt {
				// Send message to panes to resize themselves to expand back
				// into space occupied by prompt.
				m.mode = normalMode
				_ = m.PaneManager.Update(tea.WindowSizeMsg{
					Height: m.viewHeight(),
					Width:  m.viewWidth(),
				})
			}
			return m, cmd
		}

		switch {
		case key.Matches(msg, keys.Global.Quit):
			return m, tui.YesNoPrompt("Quit?", tea.Quit)
		case key.Matches(msg, keys.Global.Help):
			// '?' toggles help widget
			m.showHelp = !m.showHelp
			// Help widget takes up space so update panes' dimensions
			m.PaneManager.Update(tea.WindowSizeMsg{
				Height: m.viewHeight(),
				Width:  m.viewWidth(),
			})

		// case key.Matches(msg, keys.Global.Logs):
		// 	return m, tui.NavigateTo(tui.LogListKind)

		default:
			// Send all other keys to panes.
			if cmd := m.PaneManager.Update(msg); cmd != nil {
				return m, cmd
			}
			// If pane manager doesn't respond with a command, then send key to
			// any updateable model makers; first one to respond with a command
			// wins.
			// for _, maker := range m.makers {
			// 	if updateable, ok := maker.(updateableMaker); ok {
			// 		if cmd := updateable.Update(msg); cmd != nil {
			// 			return m, cmd
			// 		}
			// 	}
			// }

			return m, nil
		}
	case tui.ErrorMsg:
		m.err = error(msg)
	case tui.InfoMsg:
		m.info = string(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// log.Printf("top.model size: w:%d, h:%d", msg.Width, msg.Height)

		m.PaneManager.Update(tea.WindowSizeMsg{
			Height: m.viewHeight(),
			Width:  m.viewWidth(),
		})
	case cursor.BlinkMsg:
		// Send blink message to prompt if in prompt mode otherwise forward it
		// to the active pane to handle.
		if m.mode == promptMode {
			cmd = m.prompt.HandleBlink(msg)
		} else {
			// cmd = m.FocusedModel().Update(msg)
		}
		return m, cmd
	default:
		// Send remaining msg types to pane manager to route accordingly.
		cmds = append(cmds, m.PaneManager.Update(msg))
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// Start composing vertical stack of components that fill entire terminal.
	var components []string

	// Add prompt if in prompt mode.
	if m.mode == promptMode {
		components = append(components, m.prompt.View(m.width))
	}

	// Add panes
	components = append(components, lipgloss.NewStyle().
		Height(m.viewHeight()).
		Width(m.viewWidth()).
		Render(m.PaneManager.View()),
	)

	// Add help if enabled
	if m.showHelp {
		components = append(components, m.help())
	}

	// Compose footer
	footer := helpWidget
	if m.err != nil {
		footer += styles.Regular.Padding(0, 1).
			Background(styles.Red).
			Foreground(styles.White).
			Width(m.availableFooterMsgWidth()).
			Render(m.err.Error())
	} else if m.info != "" {
		footer += styles.Padded.
			Foreground(styles.Black).
			Background(styles.LightGreen).
			Width(m.availableFooterMsgWidth()).
			Render(m.info)
	} else {
		footer += styles.Padded.
			Foreground(styles.Black).
			Background(styles.EvenLighterGrey).
			Width(m.availableFooterMsgWidth()).
			Render(m.info)
	}
	footer += versionWidget

	// Add footer
	components = append(components, styles.Regular.
		Inline(true).
		MaxWidth(m.width).
		Width(m.width).
		Render(footer),
	)
	return strings.Join(components, "\n")
}

var (
	helpWidget    = styles.Padded.Background(styles.Grey).Foreground(styles.White).Render(fmt.Sprintf("%s for help", keys.Global.Help.Help().Key))
	versionWidget = styles.Padded.Background(styles.DarkGrey).Foreground(styles.White).Render("0.0.1 (21.12.24)")
)

func (m model) availableFooterMsgWidth() int {
	// -2 to accommodate padding
	return max(0, m.width-lipgloss.Width(helpWidget)-lipgloss.Width(versionWidget))

	//return m.width
}

func (m model) viewHeight() int {
	vh := m.height - tui.FooterHeight
	if m.mode == promptMode {
		vh -= tui.PromptHeight
	}
	if m.showHelp {
		vh -= tui.HelpWidgetHeight
	}

	//return max(tui.MinContentHeight, vh)

	return vh
}

// viewWidth retrieves the width available within the main view
//
// TODO: rename contentWidth
func (m model) viewWidth() int {
	return max(tui.MinContentWidth, m.width)
}

var (
	helpKeyStyle  = styles.Bold.Foreground(styles.HelpKey).Margin(0, 1, 0, 0)
	helpDescStyle = styles.Regular.Foreground(styles.HelpDesc)
)

// help renders key bindings
func (m model) help() string {
	// Compile list of bindings to render
	bindings := []key.Binding{keys.Global.Help, keys.Global.Quit}
	switch m.mode {
	case promptMode:
		bindings = append(bindings, m.prompt.HelpBindings()...)
	default:
		bindings = append(bindings, m.PaneManager.HelpBindings()...)
	}
	bindings = append(bindings, keys.KeyMapToSlice(keys.Global)...)
	//bindings = append(bindings, keys.KeyMapToSlice(keys.Navigation)...)
	bindings = removeDuplicateBindings(bindings)

	// Enumerate through each group of bindings, populating a series of
	// pairs of columns, one for keys, one for descriptions
	var (
		pairs []string
		width int
		// Subtract 2 to accommodate borders
		rows = tui.HelpWidgetHeight - 2
	)
	for i := 0; i < len(bindings); i += rows {
		var (
			keys  []string
			descs []string
		)
		for j := i; j < min(i+rows, len(bindings)); j++ {
			keys = append(keys, helpKeyStyle.Render(bindings[j].Help().Key))
			descs = append(descs, helpDescStyle.Render(bindings[j].Help().Desc))
		}
		// Render pair of columns; beyond the first pair, render a three space
		// left margin, in order to visually separate the pairs.
		var cols []string
		if len(pairs) > 0 {
			cols = []string{"   "}
		}
		cols = append(cols,
			strings.Join(keys, "\n"),
			strings.Join(descs, "\n"),
		)

		pair := lipgloss.JoinHorizontal(lipgloss.Top, cols...)
		// check whether it exceeds the maximum width avail (the width of the
		// terminal, subtracting 2 for the borders).
		width += lipgloss.Width(pair)
		if width > m.width-2 {
			break
		}
		pairs = append(pairs, pair)
	}
	// Join pairs of columns and enclose in a border
	content := lipgloss.JoinHorizontal(lipgloss.Top, pairs...)
	return styles.Border.Height(rows).Width(m.width - 2).Render(content)
}

// removeDuplicateBindings removes duplicate bindings from a list of bindings. A
// binding is deemed a duplicate if another binding has the same list of keys.
func removeDuplicateBindings(bindings []key.Binding) []key.Binding {
	seen := make(map[string]struct{})
	var i int
	for _, b := range bindings {
		key := strings.Join(b.Keys(), " ")
		if _, ok := seen[key]; ok {
			// duplicate, skip
			continue
		}
		seen[key] = struct{}{}
		bindings[i] = b
		i++
	}
	return bindings[:i]
}
