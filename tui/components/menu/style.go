package menu

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7f7f7f"))

	optionsSectionContainerStyle = lipgloss.NewStyle().Padding(1, 0)

	defaultOptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7f7f7f"))

	selectedOptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f6359a")).
				Italic(true)
)
