package mnemonicInput

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7f7f7f"))

	mnemonicSectionContainerStyle = lipgloss.NewStyle().Padding(1, 0)

	inputStyle = lipgloss.NewStyle().AlignVertical(lipgloss.Left)

	defaultInputColor = lipgloss.Color("#7f7f7f")

	focusedInputColor = lipgloss.Color("#f6359a")
)
