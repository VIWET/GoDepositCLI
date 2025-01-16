package bls

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/viwet/GoDepositCLI/tui"
)

type style struct {
	title     lipgloss.Style
	container lipgloss.Style

	colors tui.Colorscheme
}

func newStyle(colors tui.Colorscheme) style {
	return style{
		title:     lipgloss.NewStyle().Bold(true),
		container: lipgloss.NewStyle().Padding(1, 0),

		colors: colors,
	}
}
