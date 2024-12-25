package password

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7f7f7f"))

	passwordSectionContainerStyle = lipgloss.NewStyle().Padding(1, 0)

	inputFormStyle = lipgloss.NewStyle().
			Foreground(defaultInputColor).
			AlignHorizontal(lipgloss.Left).
			Width(14)

	inputStyle = lipgloss.NewStyle().AlignVertical(lipgloss.Left)

	defaultInputColor = lipgloss.Color("#7f7f7f")

	focusedInputColor = lipgloss.Color("#f6359a")

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DB2842")).Padding(0, 0, 0, 2)
)
