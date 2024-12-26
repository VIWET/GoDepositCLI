package mnemonicInput

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/viwet/GoDepositCLI/tui"
)

type Model struct {
	focused int
	input   []textinput.Model

	binding bindings
	help    help.Model
}

func New() *Model {
	return &Model{
		input:   []textinput.Model{newInput()},
		binding: newBindings(),
		help:    help.New(),
	}
}

func newInput() textinput.Model {
	input := textinput.New()
	input.EchoMode = textinput.EchoPassword
	input.Prompt = ""
	return input
}

func (m *Model) Mnemonic() string {
	return ""
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.input[m.focused].Focused() {
		return m, m.input[m.focused].Focus()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.binding.accept):
			return m, tui.Quit()

		case key.Matches(msg, m.binding.toggle):
			for i := range len(m.input) {
				m.input[i].EchoMode ^= 1
			}

			return m, nil

		case key.Matches(msg, m.binding.next) && m.focused < len(m.input)-1:
			if m.input[m.focused].Value() == "" {
				return m, nil
			}

			m.input[m.focused].Blur()
			m.focused++
			return m, m.input[m.focused].Focus()

		case key.Matches(msg, m.binding.prev) && m.focused > 0:
			if m.input[m.focused].Value() == "" {
				return m, nil
			}

			m.input[m.focused].Blur()
			m.focused--
			return m, m.input[m.focused].Focus()

		case key.Matches(msg, m.binding.space):
			if m.input[m.focused].Value() == "" {
				return m, nil
			}

			if m.focused == len(m.input)-1 {
				input := newInput()
				input.EchoMode = m.input[len(m.input)-1].EchoMode
				m.input = append(m.input, input)
			}

			m.input[m.focused].Blur()
			m.focused++
			return m, m.input[m.focused].Focus()

		case key.Matches(msg, m.binding.backspace) && m.focused > 0 && m.focused == len(m.input)-1 && m.input[m.focused].Value() == "":
			m.input[m.focused].Blur()
			m.input = m.input[:m.focused]
			m.focused--
			return m, m.input[m.focused].Focus()

		case key.Matches(msg, m.binding.quit):
			return m, tui.QuitWithError(errors.New("mnemoinc input canceled"))
		}
	}

	return m, m.updateInput(msg)
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		renderTitle("Mnemonic"+" "+strconv.Itoa(m.focused)),
		m.renderForm(),
		renderHelp(m.help, m.binding),
	)
}

func (m *Model) updateInput(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	for i := range len(m.input) {
		input, cmd := m.input[i].Update(msg)
		m.input[i] = input
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func renderTitle(title string) string {
	return titleStyle.Render(title)
}

func (m *Model) renderForm() string {
	return mnemonicSectionContainerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			renderInput(m.input, m.focused),
		),
	)
}

func renderInput(input []textinput.Model, focused int) string {
	inputStyle = inputStyle.Foreground(defaultInputColor)
	views := make([]string, len(input))
	for i := range input {
		if i == focused {
			views[i] = inputStyle.Foreground(focusedInputColor).Render(input[i].View())
		} else {
			views[i] = inputStyle.Render(input[i].View())
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		views...,
	)
}

func renderHelp(help help.Model, binding help.KeyMap) string {
	return help.View(binding)
}
