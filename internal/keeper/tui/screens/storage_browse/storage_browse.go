package storagebrowse

import (
	"context"
	"fmt"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/pkg/models"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	tableBorderSize = 4
)

type StorageBrowseScreen struct {
	storage storage.Storage
	table   table.Model
	style   lipgloss.Style
}

func (s StorageBrowseScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewStorageBrowseScreenScreen(msg.Storage), nil
}

func NewStorageBrowseScreenScreen(strg storage.Storage) *StorageBrowseScreen {
	scr := &StorageBrowseScreen{
		storage: strg,
		table:   prepareTable(),
		style:   lipgloss.NewStyle(),
	}

	scr.updateRows()

	return scr
}

func (s *StorageBrowseScreen) updateRows() {
	secrets, _ := s.storage.GetAll(context.Background())

	sortSecrets(secrets)

	rows := []table.Row{}
	for _, sec := range secrets {
		rows = append(rows, table.Row{
			strconv.Itoa(int(sec.ID)),
			sec.Title,
			sec.SecretType,
			sec.CreatedAt.Format("02 Jan 06 15:04"),
			sec.UpdatedAt.Format("02 Jan 06 15:04"),
		})
	}

	s.table.SetRows(rows)
}

func prepareTable() table.Model {
	columns := []table.Column{
		{Title: "id", Width: 5},
		{Title: "Title", Width: 20},
		{Title: "SecretType", Width: 20},
		{Title: "Created", Width: 20},
		{Title: "Updated", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	st := table.DefaultStyles()
	st.Header = st.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	st.Selected = st.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(st)

	return t
}

func (s StorageBrowseScreen) Init() tea.Cmd {
	return nil
}

func (s *StorageBrowseScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		s.table.SetWidth(min(msg.Width, s.colsWidth()))
		s.table.SetHeight(msg.Height - tableBorderSize)
	case tea.KeyMsg:
		switch msg.String() {

		// TODO: help? keys?
		case "a": // add
			cmd = tui.SetBodyPane(tui.SecretTypeScreen, tui.WithStorage(s.storage))
			cmds = append(cmds, cmd)
		case "d":
			secret, err := s.getSelectedSecret()
			if err != nil {
				cmds = append(cmds, tui.ReportError(fmt.Errorf("failed to load secret: %w", err)))
				break
			}

			err = s.storage.Delete(context.Background(), secret.ID)
			if err != nil {
				cmds = append(cmds, tui.ReportError(fmt.Errorf("failed to delete secret: %w", err)))
			}
		case "e", "enter":
			secret, err := s.getSelectedSecret()
			if err != nil {
				cmds = append(cmds, tui.ReportError(fmt.Errorf("failed to load secret: %w", err)))
				break
			}

			screen, err := s.getScreenForSecret(secret)
			if err != nil {
				cmds = append(cmds, tui.ReportError(err))
				break
			}

			cmd = tui.SetBodyPane(screen, tui.WithSecret(secret), tui.WithStorage(s.storage))
			cmds = append(cmds, cmd)
		}
	}

	s.updateRows()
	s.table, cmd = s.table.Update(msg)
	s.table.Focus()

	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (s StorageBrowseScreen) View() string {
	var b strings.Builder

	st := s.style.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	b.WriteString(fmt.Sprintf("Operating storage %s\n", s.storage.String()))
	b.WriteString(fmt.Sprintf("Use ↑↓ to navigate, (a)dd, (e)dit, (d)elete, (c)opy\n")) // TODO: help bindings
	b.WriteString(st.Render(s.table.View()))

	return b.String()
}

func (s StorageBrowseScreen) getSelectedSecret() (secret *models.Secret, err error) {
	row := s.table.SelectedRow()

	secret, err = s.loadSecret(row[0])
	if err != nil {
		return nil, err
	}

	return secret, err
}

func (s StorageBrowseScreen) loadSecret(rawID string) (*models.Secret, error) {
	var err error

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		return nil, err
	}

	sec, err := s.storage.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &sec, err
}

func (s StorageBrowseScreen) getScreenForSecret(secret *models.Secret) (tui.Screen, error) {
	switch secret.SecretType {
	case string(models.CredSecret):
		return tui.CredentialEditScreen, nil
	case string(models.TextSecret):
		return tui.TextEditScreen, nil
	case string(models.BlobSecret):
		return tui.BlobEditScreen, nil
	case string(models.CardSecret):
		return tui.CardEditScreen, nil
	default:
		return -1, fmt.Errorf("unknown secret type")
	}
}

func sortSecrets(secrets []models.Secret) {
	sort.Slice(secrets, func(i, j int) bool {
		return secrets[i].UpdatedAt.After(secrets[j].UpdatedAt)
	})
}

func (s StorageBrowseScreen) colsWidth() int {
	cols := s.table.Columns()
	total := tableBorderSize
	for _, c := range cols {
		total += c.Width
	}

	return total
}
