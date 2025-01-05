package storagebrowse

import (
	"context"
	"fmt"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"
	"gophkeeper/internal/keeper/tui/styles"
	"gophkeeper/pkg/models"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	tableBorderSize = 4
)

type savePathMsg = struct {
	path   string
	secret *models.Secret
}

type StorageBrowseScreen struct {
	storage storage.Storage
	table   table.Model
}

func (s StorageBrowseScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewStorageBrowseScreenScreen(msg.Storage), nil
}

func NewStorageBrowseScreenScreen(strg storage.Storage) *StorageBrowseScreen {
	scr := &StorageBrowseScreen{
		storage: strg,
		table:   prepareTable(),
	}

	scr.updateRows()

	return scr
}

func (s StorageBrowseScreen) Init() tea.Cmd {
	s.updateRows()
	return nil
}

func (s *StorageBrowseScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case savePathMsg: // msg from prompt for blob-secret copy-hotkey
		os.WriteFile(msg.path, msg.secret.Blob.FileBytes, 0644)
		cmds = append(cmds, infoCmd("file saved successfully"))
	case tea.WindowSizeMsg:
		s.table.SetWidth(min(msg.Width, s.colsWidth()))
		s.table.SetHeight(msg.Height - tableBorderSize)
	case tea.KeyMsg:
		switch msg.String() {
		case "a": // add
			cmd = tui.SetBodyPane(tui.SecretTypeScreen, tui.WithStorage(s.storage))
			cmds = append(cmds, cmd)
		case "e", "enter": // edit
			cmds = append(cmds, s.handleEdit())
		case "c":
			cmds = append(cmds, s.handleCopy())
		case "d": // delete
			cmds = append(cmds, s.handleDelete())

			// update table
			s.updateRows()
			cmds = append(cmds, tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(s.storage)))
		}
	}

	s.table.Focus()
	s.table, cmd = s.table.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (s StorageBrowseScreen) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Operating storage %s\n", styles.Highlighted.Render(s.storage.String())))
	b.WriteString(fmt.Sprintf("Use ↑↓ to navigate, (a)dd, (e)dit, (d)elete, (c)opy\n"))
	b.WriteString(tableStyle.Render(s.table.View()))

	return screenStyle.Render(b.String())
}

func (s *StorageBrowseScreen) HelpBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add secret")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit secret")),
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete secret")),
		key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy/save secret")),
	}
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

func (s StorageBrowseScreen) handleEdit() tea.Cmd {
	secret, err := s.getSelectedSecret()
	if err != nil {
		return errCmd("failed to load secret: %w", err)
	}

	screen, err := s.getScreenForSecret(secret)
	if err != nil {
		return errCmd("failed to get screen: %w", err)
	}

	return tui.SetBodyPane(screen, tui.WithSecret(secret), tui.WithStorage(s.storage))
}

func (s StorageBrowseScreen) handleCopy() tea.Cmd {
	secret, err := s.getSelectedSecret()
	if err != nil {
		return errCmd("failed to load secret: %w", err)
	}

	if secret.SecretType == string(models.BlobSecret) {
		// prompt and save file
		return tui.StringPrompt("choose path to save", func(str string) tea.Cmd { return func() tea.Msg { return savePathMsg{path: str, secret: secret} } })
	}

	if err := clipboard.WriteAll(secret.ToClipboard()); err != nil {
		return errCmd("failed to copy to clipboard: %w", err)
	}

	return infoCmd("secret copied successfully")
}

func (s StorageBrowseScreen) handleDelete() tea.Cmd {
	secret, err := s.getSelectedSecret()
	if err != nil {
		return errCmd("failed to load secret", err)
	}

	err = s.storage.Delete(context.Background(), secret.ID)
	if err != nil {
		return errCmd("failed to delete secret", err)
	}

	return infoCmd("secret deleted")
}

func errCmd(msg string, err error) tea.Cmd {
	return tui.ReportError(fmt.Errorf("%s: %w", msg, err))
}

func infoCmd(msg string) tea.Cmd {
	return tui.ReportInfo("%s", msg)
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

	return sec, err
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

func (s StorageBrowseScreen) colsWidth() int {
	cols := s.table.Columns()
	total := tableBorderSize
	for _, c := range cols {
		total += c.Width
	}

	return total
}

func sortSecrets(secrets []*models.Secret) {
	sort.Slice(secrets, func(i, j int) bool {
		return secrets[i].UpdatedAt.After(secrets[j].UpdatedAt) // UpdatedAt desc
	})
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
	st.Header = tableHeaderStyle
	st.Selected = tableSelectedStyle
	t.SetStyles(st)

	return t
}
