package menu

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	title string

	options []Option
	focused int

	binding bindings
	help    help.Model
}

func New(title string, options ...Option) tea.Model {
	return Model{
		title:   title,
		options: options,
		binding: newBindings(),
		help:    help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.binding.up):
			m.focused = max(m.focused-1, 0)
		case key.Matches(msg, m.binding.down):
			m.focused = min(m.focused+1, len(m.options)-1)
		case key.Matches(msg, m.binding.accept):
			model, cmd := m.options[m.focused].action()
			if model != nil {
				return model, cmd
			}
			return m, cmd
		case key.Matches(msg, m.binding.quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		renderTitle(m.title),
		renderOptions(m.options, m.focused),
		renderHelp(m.help, m.binding),
	)
}

func renderTitle(title string) string {
	return titleStyle.Render(title)
}

func renderOptions(options []Option, focused int) string {
	opts := make([]string, len(options))
	for i, opt := range options {
		if i == focused {
			opts[i] = renderSelected(opt)
		} else {
			opts[i] = renderDefault(opt)
		}
	}

	return optionsSectionContainerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			opts...,
		),
	)
}

func renderHelp(help help.Model, binding help.KeyMap) string {
	return help.View(binding)
}
