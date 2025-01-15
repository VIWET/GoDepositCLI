package menu

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/viwet/GoDepositCLI/tui"
)

type style struct {
	title     lipgloss.Style
	container lipgloss.Style

	selected lipgloss.Style
	option   lipgloss.Style

	colors tui.Colorscheme
}

func newStyle(colors tui.Colorscheme) style {
	return style{
		title:     lipgloss.NewStyle().Bold(true),
		container: lipgloss.NewStyle().Padding(1, 0),
		selected:  lipgloss.NewStyle().Italic(true),
		option:    lipgloss.NewStyle(),
		colors:    colors,
	}
}
