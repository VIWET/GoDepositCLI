package mnemonic_input

import (
	tea "github.com/charmbracelet/bubbletea"
)

type mnemonic struct {
	mnemonic []string
}

func Confirm(m []string) tea.Cmd {
	return func() tea.Msg {
		return mnemonic{m}
	}
}

type next struct {
	prev int
}

func Next(curr int) tea.Cmd {
	return func() tea.Msg {
		return next{curr}
	}
}

type prev struct {
	next int
}

func Prev(curr int) tea.Cmd {
	return func() tea.Msg {
		return prev{curr}
	}
}

type errorMsg struct{ err error }

func Error(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg{err}
	}
}
