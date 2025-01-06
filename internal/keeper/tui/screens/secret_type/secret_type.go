package secrettype

import (
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/screens"
	"gophkeeper/pkg/models"
	"maps"
	"slices"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	selectBack = iota
	selectCredential
	selectText
	selectCard
	selectBlob
)

// Model which renders selection list of different commands
type SecretTypeScreen struct {
	tea.Model

	storage storage.Storage

	list   list.Model
	choice string
}

func (s SecretTypeScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewSecretTypeScreen(msg.Storage), nil
}

func NewSecretTypeScreen(strg storage.Storage) *SecretTypeScreen {
	m := &SecretTypeScreen{storage: strg}
	m.prepareSecretListModel()

	return m
}

func (s *SecretTypeScreen) prepareSecretListModel() {
	choices := map[int]string{
		selectBack:       "Go back",
		selectCredential: "Add credentials",
		selectText:       "Add text",
		selectCard:       "Add card info",
		selectBlob:       "Upload file",
	}

	keys := slices.Collect(maps.Keys(choices))
	sort.Ints(keys)

	items := []list.Item{}
	for i := range keys {
		items = append(items, secretItem{id: i, name: choices[i]})
	}

	l := list.New(items, secretItemDelegate{}, 0, 0)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.KeyMap.Quit.SetEnabled(false)

	s.list = l
}

func (s SecretTypeScreen) Init() tea.Cmd {
	return tea.SetWindowTitle("GophKeeper client")
}

func (s *SecretTypeScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // pane size
		s.list.SetWidth(msg.Width)
		s.list.SetHeight(msg.Height - 2 - 2)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":

			i, ok := s.list.SelectedItem().(secretItem)
			if ok {
				s.choice = string(i.name)
				// cmds = append(cmds, s.items[string(i)].cmd)
			}

			switch i.id {
			case selectBack:
				cmd = tui.SetBodyPane(
					tui.StorageBrowseScreen,
					tui.WithStorage(s.storage),
				)

				return cmd

			case selectCredential:
				sec := models.NewSecret(models.CredSecret)

				cmd = tui.SetBodyPane(
					tui.CredentialEditScreen,
					tui.WithStorage(s.storage),
					tui.WithSecret(sec),
				)
			case selectText:
				sec := models.NewSecret(models.TextSecret)

				cmd = tui.SetBodyPane(
					tui.TextEditScreen,
					tui.WithStorage(s.storage),
					tui.WithSecret(sec),
				)
			case selectCard:
				sec := models.NewSecret(models.CardSecret)

				cmd = tui.SetBodyPane(
					tui.CardEditScreen,
					tui.WithStorage(s.storage),
					tui.WithSecret(sec),
				)
			case selectBlob:
				sec := models.NewSecret(models.BlobSecret)

				cmd = tui.SetBodyPane(
					tui.BlobEditScreen,
					tui.WithStorage(s.storage),
					tui.WithSecret(sec),
				)
			}

			cmds = append(cmds, cmd)
		}
	}

	s.list, cmd = s.list.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (s *SecretTypeScreen) View() string {
	return screens.RenderContent("Select type of secret:", s.list.View())
}
