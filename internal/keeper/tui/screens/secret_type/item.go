package secrettype

import (
	"fmt"
	"gophkeeper/internal/keeper/tui/styles"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	var (
		itemStyle         = styles.Regular.PaddingLeft(4)
		itemSelectedStyle = styles.Regular.PaddingLeft(2).Foreground(lipgloss.Color("170"))
	)

	i, ok := listItem.(secretItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return itemSelectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
