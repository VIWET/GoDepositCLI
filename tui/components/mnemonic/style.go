package mnemonic

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7f7f7f"))

	mnemonicSectionContainerStyle = lipgloss.NewStyle().Padding(1, 0)

	mnemonicLanguageSectionContainerStyle = lipgloss.NewStyle().Padding(0, 0, 1, 0)

	mnemonicLanguageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7f7f7f"))

	mnemonicWordStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7f7f7f"))

	mnemonicIndexStyle = lipgloss.NewStyle().
				Width(3).
				Foreground(lipgloss.Color("#4f4f4f"))

	mnemonicWordIndexStyle = lipgloss.NewStyle().
				Width(20).
				AlignHorizontal(lipgloss.Left)
)
