package tui

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const logFilePath = "debug.log"

var ExitCmdErr = errors.New("exit command")

// App describes TUI app
type App struct {
	startAt time.Time

	windowHeight int
	windowWidth  int

	notify chan error

	logFile *os.File

	//root *screens.RootModel
}

func NewApp() *App {
	return &App{
		notify: make(chan error, 1),
	}
}

func (a App) Start() {

	var err error

	a.setupLogging()

	//state := screens.NewState()
	// root := screens.NewWindow()
	// root := screens.GetWindow()

	// p := tea.NewProgram(root, tea.WithAltScreen())

	// Run tea program
	go func() {

		//	err = top.Start()
		//_, err = p.Run()
		if err != nil {
			log.Fatal("failed to run bubbletea app: ", err)
		}

		a.notify <- ExitCmdErr
	}()

}

func (a App) Shutdown() {
	a.logFile.Close()
}

func (a App) Notify() chan error {
	return a.notify
}

func (a *App) SetSize(w, h int) {
	a.windowHeight = h
	a.windowWidth = w
}

func (a App) Height() int {
	return a.windowHeight
}

func (a App) Width() int {
	return a.windowWidth
}

func (a *App) setupLogging() {
	var err error
	a.logFile, err = tea.LogToFile(logFilePath, "debug")
	if err != nil {
		a.notify <- fmt.Errorf("failed to setup bubbletea log: %w", err)
		return
	}
}
