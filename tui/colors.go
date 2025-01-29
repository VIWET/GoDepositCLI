package tui

import "github.com/charmbracelet/lipgloss"

type Colorscheme struct {
	Text   lipgloss.AdaptiveColor
	Title  lipgloss.AdaptiveColor
	Accent lipgloss.AdaptiveColor
	Error  lipgloss.AdaptiveColor
}

func DefaultColorscheme() Colorscheme {
	return Colorscheme{
		Text:   lipgloss.AdaptiveColor{Light: "#3B4252", Dark: "#7b8b99"},
		Title:  lipgloss.AdaptiveColor{Light: "#E5E9F0", Dark: "#dcd7d7"},
		Accent: lipgloss.AdaptiveColor{Light: "#B48EAD", Dark: "#bc7fd2"},
		Error:  lipgloss.AdaptiveColor{Light: "#BF616A", Dark: "#df2683"},
	}
}
