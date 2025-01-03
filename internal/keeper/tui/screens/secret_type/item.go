package secrettype

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	stitemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	stitemSelectedStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type secretItem struct {
	id   int
	name string
}

func (i secretItem) FilterValue() string { return "" }

type secretItemDelegate struct{}

func (d secretItemDelegate) Height() int                             { return 1 }
func (d secretItemDelegate) Spacing() int                            { return 0 }
func (d secretItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d secretItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(secretItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.name)

	fn := stitemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return stitemSelectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
