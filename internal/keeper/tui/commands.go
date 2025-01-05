package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ErrorMsg error
type InfoMsg string

func ReportInfo(msg string, args ...any) tea.Cmd {
	return CmdHandler(InfoMsg(fmt.Sprintf(msg, args...)))
}

func ReportError(err error) tea.Cmd {
	return CmdHandler(ErrorMsg(err))
}

func NavigateTo(screen Screen, opts ...NavigateOption) tea.Cmd {
	return CmdHandler(NewNavigationMsg(screen, opts...))
}

func SetBodyPane(screen Screen, opts ...NavigateOption) tea.Cmd {
	opts = append(opts, WithPosition(BodyPane))
	return NavigateTo(screen, opts...)
}

func GoToStart() tea.Cmd {
	return tea.Batch(
		NavigateTo(WelcomeScreen, []NavigateOption{DisableFocus(), WithPosition(BodyPane)}...),
		NavigateTo(MenuScreen, []NavigateOption{WithPosition(LeftPane)}...),
	)
}

func SetLeftPane(screen Screen) tea.Cmd {
	return NavigateTo(screen, WithPosition(LeftPane))
}
