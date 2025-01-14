package tui

import "github.com/charmbracelet/lipgloss"

type Colorscheme struct {
	Black   lipgloss.AdaptiveColor
	Red     lipgloss.AdaptiveColor
	Green   lipgloss.AdaptiveColor
	Yellow  lipgloss.AdaptiveColor
	Blue    lipgloss.AdaptiveColor
	Magenta lipgloss.AdaptiveColor
	Cyan    lipgloss.AdaptiveColor
	White   lipgloss.AdaptiveColor
}

func DefaultColorscheme() Colorscheme {
	return Colorscheme{
		Black:   lipgloss.AdaptiveColor{Light: "#3B4252", Dark: "#7b8b99"},
		Red:     lipgloss.AdaptiveColor{Light: "#BF616A", Dark: "#df2683"},
		Green:   lipgloss.AdaptiveColor{Light: "#A3BE8C", Dark: "#13868c"},
		Yellow:  lipgloss.AdaptiveColor{Light: "#EBCB8B", Dark: "#fcfcdf"},
		Blue:    lipgloss.AdaptiveColor{Light: "#81A1C1", Dark: "#1a86b9"},
		Magenta: lipgloss.AdaptiveColor{Light: "#B48EAD", Dark: "#bc7fd2"},
		Cyan:    lipgloss.AdaptiveColor{Light: "#88C0D0", Dark: "#7cc7d6"},
		White:   lipgloss.AdaptiveColor{Light: "#E5E9F0", Dark: "#dcd7d7"},
	}
}
