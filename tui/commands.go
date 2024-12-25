package tui

import tea "github.com/charmbracelet/bubbletea"

type quit struct {
	err error
}

func Quit() tea.Cmd {
	return func() tea.Msg {
		return quit{}
	}
}

func QuitWithError(err error) tea.Cmd {
	return func() tea.Msg {
		return quit{err: err}
	}
}
