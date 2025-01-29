package password

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/viwet/GoDepositCLI/tui"
)

type style struct {
	title     lipgloss.Style
	container lipgloss.Style

	form  lipgloss.Style
	input lipgloss.Style
	error lipgloss.Style

	colors tui.Colorscheme
}

func newStyle(colors tui.Colorscheme) style {
	return style{
		title:     lipgloss.NewStyle().Bold(true),
		container: lipgloss.NewStyle().Padding(1, 0),

		form:  lipgloss.NewStyle().Width(14).AlignHorizontal(lipgloss.Left),
		input: lipgloss.NewStyle(),
		error: lipgloss.NewStyle().Padding(0, 0, 0, 2),

		colors: colors,
	}
}
