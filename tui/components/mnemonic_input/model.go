package mnemonicInput

import (
	"errors"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/viwet/GoDepositCLI/tui"
)

const columns = 3

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
	input.Width = 15
	input.CharLimit = 15
	input.KeyMap = inputBinding
	input.Validate = func(value string) error {
		if strings.Contains(value, " ") {
			return errors.New("paste is not allowed")
		}
		return nil
	}
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
		renderTitle("Mnemonic"),
		m.renderForm(),
		renderHelp(m.help, m.binding),
	)
}

func (m *Model) updateInput(msg tea.Msg) tea.Cmd {
	var (
		cmds   []tea.Cmd
		inputs = m.input
		offset = 0
	)

	for i := range len(m.input) {
		input, cmd := inputs[i].Update(msg)
		if input.Err != nil {
			m.input[m.focused].Blur()

			words := strings.Fields(input.Value())
			newInputs := make([]textinput.Model, len(m.input)-1+len(words))

			// Copy inputs
			copy(newInputs[:i], inputs[:i])
			for j := 0; j < len(words); j++ {
				newInput := newInput()
				newInput.SetValue(words[j])
				if len(inputs) > 0 {
					newInput.EchoMode = inputs[0].EchoMode
				}
				newInputs[i+j] = newInput
				offset++
			}
			copy(newInputs[i+len(words):], inputs[i+1:])
			inputs = newInputs

			m.focused = i + len(words) - 1
		} else {
			inputs[i] = input
		}
		cmds = append(cmds, cmd)
	}
	m.input = inputs
	return tea.Batch(cmds...)
}

func renderTitle(title string) string {
	return titleStyle.Render(title)
}

func (m *Model) renderForm() string {
	return mnemonicSectionContainerStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			renderInput(m.input, 0, m.focused),
			renderInput(m.input, 1, m.focused),
			renderInput(m.input, 2, m.focused),
		),
	)
}

func renderInput(input []textinput.Model, column, focused int) string {
	var (
		rows       = len(input)/columns + 1
		views      = make([]string, 0, rows)
		inputStyle = inputStyle.Foreground(defaultInputColor)
	)

	for row := range rows {
		index := row*columns + column

		if index >= len(input) {
			break
		}

		var word string
		if index == focused {
			word = renderWordWithIndex(inputStyle.Foreground(focusedInputColor).Render(input[index].View()), index+1)
		} else {
			word = renderWordWithIndex(inputStyle.Render(input[index].View()), index+1)
		}
		views = append(views, word)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		views...,
	)
}

func renderWordWithIndex(word string, index int) string {
	return mnemonicWordIndexStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			mnemonicIndexStyle.Render(strconv.Itoa(index)),
			word,
		),
	)
}

func renderHelp(help help.Model, binding help.KeyMap) string {
	return help.View(binding)
}
