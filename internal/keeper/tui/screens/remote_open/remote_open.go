package remoteeopen

import (
	"gophkeeper/internal/keeper/api"
	"gophkeeper/internal/keeper/crypto"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/internal/keeper/tui"

	tea "github.com/charmbracelet/bubbletea"
)

type RemoteOpenScreen struct {
	client api.IApiClient
}

type RemoteOpenScreenMaker struct {
	Client api.IApiClient
}

func (m RemoteOpenScreenMaker) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewRemoteOpenScreen(m.Client), nil
}

func (s RemoteOpenScreen) Make(msg tui.NavigationMsg, width, height int) (tui.Teable, error) {
	return NewRemoteOpenScreen(msg.Client), nil
}

func NewRemoteOpenScreen(client api.IApiClient) *RemoteOpenScreen {
	return &RemoteOpenScreen{
		client: client,
	}
}

func (s RemoteOpenScreen) Init() tea.Cmd {
	var cmds []tea.Cmd

	if len(s.client.GetToken()) > 0 {
		// already authorized
		encrypter := crypto.NewKeeperEncrypter()
		strg, err := storage.NewRemoteStorage(s.client, encrypter)
		if err != nil {
			cmds = append(cmds, tui.ReportError(err))
		} else {
			cmds = append(cmds, tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(strg)))
		}
	} else {
		// go to login
		cmds = append(cmds, tui.SetBodyPane(tui.LoginScreen, tui.WithClient(s.client)))
	}

	return tea.Batch(cmds...)
}

func (s *RemoteOpenScreen) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (s RemoteOpenScreen) View() string {
	return ""
}
