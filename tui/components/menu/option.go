package menu

import tea "github.com/charmbracelet/bubbletea"

type Option struct {
	title  string
	action func() (tea.Model, tea.Cmd)
}

func NewOption(title string, action func() (tea.Model, tea.Cmd)) Option {
	return Option{
		title:  title,
		action: action,
	}
}

func renderSelected(option Option) string {
	return selectedOptionStyle.Render(option.title)
}

func renderDefault(option Option) string {
	return defaultOptionStyle.Render(option.title)
}
