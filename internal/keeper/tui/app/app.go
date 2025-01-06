package app

import (
	"errors"
	"fmt"
	"gophkeeper/internal/keeper/api"
	"gophkeeper/internal/keeper/config"
	"gophkeeper/internal/keeper/tui/top"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/dig"
)

const logFilePath = "debug.log"

var ErrExitCmd = errors.New("exit command")

// App describes TUI app
type App struct {
	config   *config.Config
	client   api.IApiClient
	deps     AppDependencies
	topModel *top.Model

	logFile *os.File
	notify  chan error
}

type AppDependencies struct {
	dig.In

	TopModel *top.Model
	Config   *config.Config
	Client   api.IApiClient
}

func NewApp(deps AppDependencies) *App {
	return &App{
		topModel: deps.TopModel,
		config:   deps.Config,
		client:   deps.Client,
		notify:   make(chan error, 1),
		deps:     deps,
	}
}

func (a App) Start() {
	a.setupLogging()

	p := tea.NewProgram(a.topModel, tea.WithAltScreen())

	// Run tea program
	go func() {
		_, err := p.Run()
		if err != nil {
			log.Fatal("failed to run bubbletea app: ", err)
		}

		a.notify <- ErrExitCmd
	}()
}

func (a App) Shutdown() {
	err := a.logFile.Close()
	if err != nil {
		panic(err)
	}
}

func (a App) Notify() chan error {
	return a.notify
}

func (a *App) setupLogging() {
	var err error

	a.logFile, err = tea.LogToFile(logFilePath, "debug")
	if err != nil {
		a.notify <- fmt.Errorf("failed to setup bubbletea log: %w", err)
		return
	}
}
