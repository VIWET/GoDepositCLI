package menu

import (
	"errors"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/viwet/GoDepositCLI/tui"
)

type Model struct {
	title string

	options []Option
	focused int

	bindings bindings
	style    style

	help help.Model
}

func New(title string, options ...Option) Model {
	return Model{
		title:    title,
		options:  options,
		bindings: newBindings(),
		style:    newStyle(tui.DefaultColorscheme()),
		help:     help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.bindings.up):
			m.focused = max(m.focused-1, 0)
		case key.Matches(msg, m.bindings.down):
			m.focused = min(m.focused+1, len(m.options)-1)
		case key.Matches(msg, m.bindings.accept):
			return m.options[m.focused].action()
		case key.Matches(msg, m.bindings.quit):
			return m, tui.QuitWithError(errors.New("selection canceled"))
		}
	}

	return m, nil
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.style.title.Foreground(m.style.colors.White).Render(m.title),
		m.style.container.Render(
			lipgloss.JoinVertical(lipgloss.Left, m.optionsView()...),
		),
		m.help.View(m.bindings),
	)
}

func (m Model) optionsView() []string {
	var (
		views    = make([]string, len(m.options))
		selected = m.style.selected.Foreground(m.style.colors.Magenta).Render
		option   = m.style.option.Foreground(m.style.colors.Black).Render
	)

	for i, opt := range m.options {
		if i == m.focused {
			views[i] = selected(opt.title)
		} else {
			views[i] = option(opt.title)
		}
	}

	return views
}
