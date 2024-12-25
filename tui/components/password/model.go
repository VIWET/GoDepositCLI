package password

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/viwet/GoDepositCLI/tui"
)

const MinPasswordLength = 8

type Model struct {
	password    textinput.Model
	passwordErr error

	confirm    textinput.Model
	confirmErr error

	binding bindings
	help    help.Model
}

func New() *Model {
	return &Model{
		password: newInput(),
		confirm:  newInput(),
		binding:  newBindings(),
		help:     help.New(),
	}
}

func newInput() textinput.Model {
	input := textinput.New()
	input.EchoMode = textinput.EchoPassword
	input.Prompt = ""
	input.CharLimit = 64
	return input
}

func (m *Model) Value() string {
	return m.password.Value()
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.password.Focused() && !m.confirm.Focused() {
		return m, m.password.Focus()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.binding.accept):
			switch {
			case m.password.Focused():
				if m.isValidPassword() {
					m.password.Blur()
					return m, m.confirm.Focus()
				}

			case m.confirm.Focused():
				if m.isConfirmedPassword() {
					return m, tui.Quit()
				}
			}

		case key.Matches(msg, m.binding.cancel):
			switch {
			case m.password.Focused():
				m.confirm.Reset()
				m.password.Reset()
				return m, nil
			case m.confirm.Focused():
				m.confirm.Reset()
				m.confirm.Blur()
				return m, m.password.Focus()
			}

		case key.Matches(msg, m.binding.toggle):
			m.password.EchoMode ^= 1
			m.confirm.EchoMode ^= 1
			return m, nil

		case key.Matches(msg, m.binding.quit):
			return m, tui.QuitWithError(errors.New("password input canceled"))
		}
	}

	return m, tea.Batch(m.updatePassword(msg), m.updateConfirm(msg))
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		renderTitle("Password"),
		m.renderForm(),
		renderHelp(m.help, m.binding),
	)
}

func (m *Model) updatePassword(msg tea.Msg) tea.Cmd {
	password, cmd := m.password.Update(msg)
	m.password = password
	return cmd
}

func (m *Model) updateConfirm(msg tea.Msg) tea.Cmd {
	confirm, cmd := m.confirm.Update(msg)
	m.confirm = confirm
	return cmd
}

func (m *Model) isValidPassword() bool {
	if strings.TrimSpace(m.password.Value()) == "" {
		m.passwordErr = errors.New("Password cannot contain only blank spaces")
		return false
	}
	if utf8.RuneCountInString(m.password.Value()) < MinPasswordLength {
		m.passwordErr = errors.New("Password must be at least 8 characters")
		return false
	}
	m.passwordErr = nil
	return true
}

func (m *Model) isConfirmedPassword() bool {
	if !m.isValidPassword() {
		return false
	}

	if m.password.Value() != m.confirm.Value() {
		m.confirmErr = errors.New("Passwords are not equal")
		return false
	}

	m.confirmErr = nil
	return true
}

func renderTitle(title string) string {
	return titleStyle.Render(title)
}

func (m *Model) renderForm() string {
	return passwordSectionContainerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			renderInput(m.password, "Password:", m.passwordErr),
			renderInput(m.confirm, "Confirmation:", m.confirmErr),
		),
	)
}

func renderInput(input textinput.Model, name string, err error) string {
	inputStyle = inputStyle.Foreground(defaultInputColor)
	if input.Focused() {
		inputStyle = inputStyle.Foreground(focusedInputColor)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		inputFormStyle.Render(name),
		inputStyle.Render(input.View()),
		renderInputError(err),
	)
}

func renderInputError(err error) string {
	if err == nil {
		return ""
	}

	return errorStyle.Render(fmt.Sprintf("[Error]: %s", err.Error()))
}

func renderHelp(help help.Model, binding help.KeyMap) string {
	return help.View(binding)
}
