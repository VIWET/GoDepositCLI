package mnemonic

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/viwet/GoDepositCLI/tui"
)

type style struct {
	title     lipgloss.Style
	container lipgloss.Style

	wordContainer lipgloss.Style
	index         lipgloss.Style
	word          lipgloss.Style

	colors tui.Colorscheme
}

func newStyle(colors tui.Colorscheme) style {
	return style{
		title:     lipgloss.NewStyle().Bold(true),
		container: lipgloss.NewStyle().Padding(1, 0),

		wordContainer: lipgloss.NewStyle().Width(20).AlignHorizontal(lipgloss.Left),
		index:         lipgloss.NewStyle().Width(3),
		word:          lipgloss.NewStyle(),

		colors: colors,
	}
}
