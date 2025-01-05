package tui

import (
	"gophkeeper/internal/keeper/api"
	"gophkeeper/internal/keeper/storage"
	"gophkeeper/pkg/models"

	tea "github.com/charmbracelet/bubbletea"
)

type ReloadSecretList struct{} // TODO

type NavigationCallback func(args ...any) tea.Cmd

// NavigationMsg is an instruction to navigate to a page
type NavigationMsg struct {
	Screen       Screen
	Page         Page
	Position     Position
	DisableFocus bool

	Callback NavigationCallback
	Client   api.IApiClient
	Secret   *models.Secret
	Storage  storage.Storage
}

func NewNavigationMsg(screen Screen, opts ...NavigateOption) NavigationMsg {
	msg := NavigationMsg{Page: Page{Screen: screen}}
	for _, fn := range opts {
		fn(&msg)
	}
	return msg
}

type NavigateOption func(msg *NavigationMsg)

func WithCallback(c NavigationCallback) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Callback = c
	}
}

func WithClient(c api.IApiClient) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Client = c
	}
}

func WithPosition(position Position) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Position = position
	}
}

func WithStorage(strg storage.Storage) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Storage = strg
	}
}

func WithSecret(sec *models.Secret) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Secret = sec
	}
}

func DisableFocus() NavigateOption {
	return func(msg *NavigationMsg) {
		msg.DisableFocus = true
	}
}
