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
	"github.com/viwet/GoDepositCLI/tui"
	"github.com/viwet/GoDepositCLI/tui/components/bls"
	"github.com/viwet/GoDepositCLI/tui/components/deposits"
)

type Model[Config app.ConfigConstraint] struct {
	ctx   *cli.Context
	state *app.State[Config]

	password textinput.Model
	confirm  textinput.Model

	bindings bindings
	style    style
	help     help.Model

	next tui.NewModel[Config]
}

func NewDepositPassword() tui.NewModel[app.DepositConfig] {
	return func(ctx *cli.Context, state *app.State[app.DepositConfig]) (tea.Model, tea.Cmd) {
		return newModel(ctx, state, deposits.New)
	}
}

func NewBLSToExecutionPassword() tui.NewModel[app.BLSToExecutionConfig] {
	return func(ctx *cli.Context, state *app.State[app.BLSToExecutionConfig]) (tea.Model, tea.Cmd) {
		return newModel(ctx, state, bls.New)
	}
}

func newModel[Config app.ConfigConstraint](ctx *cli.Context, state *app.State[Config], next tui.NewModel[Config]) (tea.Model, tea.Cmd) {
	if ctx.IsSet(tui.PasswordFlagName) {
		password := ctx.String(tui.PasswordFlagName)
		state.WithPassword(password)
		return next(ctx, state)
	}

	model := &Model[Config]{
		ctx:   ctx,
		state: state,

		password: newInput(),
		confirm:  newInput(),
		bindings: newBindings(),
		style:    newStyle(tui.DefaultColorscheme()),
		help:     help.New(),
		next:     next,
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

func (m *Model[Config]) Init() tea.Cmd {
	return nil
}

func (m *Model[Config]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *Model[Config]) View() string {
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

func (m *Model[Config]) onConfirm(msg confirm) (tea.Model, tea.Cmd) {
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
			return m.next(m.ctx, m.state)
		}
	}

	return m, nil
}

func (m *Model[Config]) onReset() (tea.Model, tea.Cmd) {
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

func (m *Model[Config]) onToggle() (tea.Model, tea.Cmd) {
	m.password.EchoMode ^= 1
	m.confirm.EchoMode ^= 1
	return m, nil
}

func (m *Model[Config]) resetInputErrors() {
	m.password.Err = nil
	m.confirm.Err = nil
}

func (m *Model[Config]) updatePassword(msg tea.Msg) tea.Cmd {
	password, cmd := m.password.Update(msg)
	m.password = password
	return cmd
}

func (m *Model[Config]) updateConfirm(msg tea.Msg) tea.Cmd {
	confirm, cmd := m.confirm.Update(msg)
	m.confirm = confirm
	return cmd
}

func (m *Model[Config]) inputView(input textinput.Model, name string) string {
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

func (m *Model[Config]) passwordView() string {
	return m.inputView(m.password, "Password:")
}

func (m *Model[Config]) confirmView() string {
	return m.inputView(m.confirm, "Confirmation:")
}
