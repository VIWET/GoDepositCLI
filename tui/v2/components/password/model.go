package password

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui/v2"
	"github.com/viwet/GoDepositCLI/tui/v2/components/deposits"
)

type Model struct {
	ctx   *cli.Context
	state *app.State[app.DepositConfig]

	password textinput.Model
	confirm  textinput.Model

	bindings bindings
	style    style
	help     help.Model
}

func New(ctx *cli.Context, state *app.State[app.DepositConfig]) (tea.Model, tea.Cmd) {
	model := &Model{
		ctx:   ctx,
		state: state,

		password: newInput(),
		confirm:  newInput(),
		bindings: newBindings(),
		style:    newStyle(tui.DefaultColorscheme()),
		help:     help.New(),
	}

	return model, model.password.Focus()
}

func newInput() textinput.Model {
	input := textinput.New()
	input.EchoMode = textinput.EchoPassword
	input.Prompt = ""
	input.CharLimit = 64
	return input
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case confirm:
		return m.onConfirm(msg)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.bindings.accept):
			return m, Confirm(m.password.Value(), m.confirm.Value())

		case key.Matches(msg, m.bindings.reset):
			return m.onReset()

		case key.Matches(msg, m.bindings.toggle):
			return m.onToggle()

		case key.Matches(msg, m.bindings.quit):
			return m, tui.QuitWithError(errors.New("password input canceled"))
		}
	}

	return m, tea.Batch(m.updatePassword(msg), m.updateConfirm(msg))
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.style.title.Foreground(m.style.colors.White).Render("Password"),
		m.style.container.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.passwordView(),
				m.confirmView(),
			),
		),
		m.help.View(m.bindings),
	)
}

func (m *Model) onConfirm(msg confirm) (tea.Model, tea.Cmd) {
	m.resetInputErrors()
	switch {
	case m.password.Focused():
		m.password.Err = msg.validation
		if msg.validation == nil {
			m.password.Blur()
			return m, m.confirm.Focus()
		}
	case m.confirm.Focused():
		m.confirm.Err = msg.confirmation
		if msg.validation == nil && msg.confirmation == nil {
			m.state.WithPassword(m.password.Value())
			return deposits.New(m.ctx, m.state)
		}
	}

	return m, nil
}

func (m *Model) onReset() (tea.Model, tea.Cmd) {
	m.resetInputErrors()
	switch {
	case m.password.Focused():
		m.confirm.Reset()
		m.password.Reset()
	case m.confirm.Focused():
		m.confirm.Reset()
		m.confirm.Blur()
	}
	return m, m.password.Focus()
}

func (m *Model) onToggle() (tea.Model, tea.Cmd) {
	m.password.EchoMode ^= 1
	m.confirm.EchoMode ^= 1
	return m, nil
}

func (m *Model) resetInputErrors() {
	m.password.Err = nil
	m.confirm.Err = nil
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

func (m *Model) inputView(input textinput.Model, name string) string {
	render := m.style.input.Foreground(m.style.colors.Black).Render
	if input.Focused() {
		render = m.style.input.Foreground(m.style.colors.Magenta).Render
	}

	var errView string
	if err := input.Err; err != nil {
		errView = m.style.error.Foreground(m.style.colors.Red).Render(
			fmt.Sprintf("[Error]: %s", err.Error()),
		)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		m.style.form.Foreground(m.style.colors.Black).Render(name),
		render(input.View()),
		errView,
	)
}

func (m *Model) passwordView() string {
	return m.inputView(m.password, "Password:")
}

func (m *Model) confirmView() string {
	return m.inputView(m.confirm, "Confirmation:")
}
