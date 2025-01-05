package top

import (
	"fmt"
	"gophkeeper/internal/keeper/api"
	"gophkeeper/internal/keeper/config"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/styles"
	"reflect"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/dig"
)

type mode int

const (
	normalMode mode = iota // default
	promptMode             // prompt is visible
)

// Top-level tea model
type Model struct {
	*tui.PaneManager

	client api.IApiClient
	makers map[tui.Screen]tui.ScreenMaker
	prompt *tui.Prompt
	mode   mode

	width    int
	height   int
	showHelp bool
	err      error
	info     string
}

type ModelDependencies struct {
	dig.In

	Config *config.Config
	Client api.IApiClient
}

func NewModel(deps ModelDependencies) (*Model, error) {
	// spinner := spinner.New(spinner.WithSpinner(spinner.Line))
	// makers := makeMakers(cfg, app, &spinner, helpers)

	makers := prepareMakers(deps)

	m := Model{
		client:      deps.Client,
		PaneManager: tui.NewPaneManager(makers),
		makers:      makers,
	}

	return &m, nil
}

func (m Model) Init() tea.Cmd {
	return m.PaneManager.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tui.PromptMsg:
		m.mode = promptMode
		var blink tea.Cmd
		m.prompt, blink = tui.NewPrompt(msg)

		// Tell panes to resize themselves
		cmd = m.PaneManager.Update(tea.WindowSizeMsg{
			Height: m.viewHeight(),
			Width:  m.viewWidth(),
		})

		return m, tea.Batch(cmd, blink)

	case tea.KeyMsg:
		m.info = "" // Clear info/error messages in the footer
		m.err = nil

		switch m.mode {
		case promptMode:
			closePrompt, cmd := m.prompt.HandleKey(msg)
			if closePrompt {
				m.mode = normalMode
				m.PaneManager.Update(tea.WindowSizeMsg{
					Height: m.viewHeight(),
					Width:  m.viewWidth(),
				})
			}
			return m, cmd
		}

		switch {
		case key.Matches(msg, tui.GlobalKeys.Quit):
			return m, tui.YesNoPrompt("Quit?", tea.Quit)
		case key.Matches(msg, tui.GlobalKeys.Help):
			m.showHelp = !m.showHelp

			m.PaneManager.Update(tea.WindowSizeMsg{
				Height: m.viewHeight(),
				Width:  m.viewWidth(),
			})

		default:
			// Send all other keys to panes
			return m, m.PaneManager.Update(msg)
		}

	case tui.ErrorMsg:
		m.err = error(msg)

	case tui.InfoMsg:
		m.info = string(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.PaneManager.Update(tea.WindowSizeMsg{
			Height: m.viewHeight(),
			Width:  m.viewWidth(),
		})

	case cursor.BlinkMsg:
		if m.mode == promptMode {
			cmd = m.prompt.HandleBlink(msg)
		} else {
			cmd = m.PaneManager.FocusedModel().Update(msg)
		}
		return m, cmd
	default:
		// Send remaining msg types to pane manager
		cmds = append(cmds, m.PaneManager.Update(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var components []string

	// Add prompt if in prompt mode
	if m.mode == promptMode {
		components = append(components, m.prompt.View(m.width))
	}

	// Add pane manager
	components = append(components, styles.Regular.
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
	helpWidget    = styles.Padded.Background(styles.Grey).Foreground(styles.White).Render(fmt.Sprintf("%s for help", tui.GlobalKeys.Help.Help().Key))
	versionWidget = styles.Padded.Background(styles.DarkGrey).Foreground(styles.White).Render("0.0.1 (21.12.24)")
)

func (m Model) availableFooterMsgWidth() int {
	return max(0, m.width-lipgloss.Width(helpWidget)-lipgloss.Width(versionWidget))
}

func (m Model) viewHeight() int {
	vh := m.height - tui.FooterHeight
	if m.mode == promptMode {
		vh -= tui.PromptHeight
	}
	if m.showHelp {
		vh -= tui.HelpWidgetHeight
	}

	// TODO: debug max(tui.MinContentHeight, vh)
	return vh
}

// Width available within the main view
func (m Model) viewWidth() int {
	return max(tui.MinContentWidth, m.width)
}

// help renders key bindings
func (m Model) help() string {
	// Compile list of bindings to render
	bindings := []key.Binding{tui.GlobalKeys.Help, tui.GlobalKeys.Quit}

	switch m.mode {
	case promptMode:
		bindings = append(bindings, m.prompt.HelpBindings()...)
	default:
		bindings = append(bindings, m.PaneManager.HelpBindings()...)
	}

	bindings = append(bindings, keyMapToSlice(tui.GlobalKeys)...)
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
			keys = append(keys, styles.HelpKeyStyle.Render(bindings[j].Help().Key))
			descs = append(descs, styles.HelpDescStyle.Render(bindings[j].Help().Desc))
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

func keyMapToSlice(t any) (bindings []key.Binding) {
	typ := reflect.TypeOf(t)
	if typ.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < typ.NumField(); i++ {
		v := reflect.ValueOf(t).Field(i)
		bindings = append(bindings, v.Interface().(key.Binding))
	}
	return
}
