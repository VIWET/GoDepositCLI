package mnemonic_input

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoBIP39/words"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui"
	"github.com/viwet/GoDepositCLI/tui/components/password"
)

const (
	mnemonicColumns  = 3
	mnemonicMinWords = 12
	mnemonicMaxWords = 24
)

type Model[Config app.ConfigConstraint] struct {
	ctx   *cli.Context
	state *app.State[Config]
	list  words.List

	echo    textinput.EchoMode
	inputs  []textinput.Model
	focused int

	style    style
	bindings bindings
	help     help.Model
	errorMsg errorMsg

	next tui.NewModel[Config]
}

func NewDepositMnemonicInput() tui.NewModel[app.DepositConfig] {
	return func(ctx *cli.Context, state *app.State[app.DepositConfig]) (tea.Model, tea.Cmd) {
		return newModel(ctx, state, state.Config().MnemonicConfig.Language, password.NewDepositPassword())
	}
}

func NewBLSToExecutionMnemonicInput() tui.NewModel[app.BLSToExecutionConfig] {
	return func(ctx *cli.Context, state *app.State[app.BLSToExecutionConfig]) (tea.Model, tea.Cmd) {
		return newModel(ctx, state, state.Config().MnemonicConfig.Language, password.NewBLSToExecutionPassword())
	}
}

func newModel[Config app.ConfigConstraint](
	ctx *cli.Context,
	state *app.State[Config],
	language string,
	next tui.NewModel[Config],
) (tea.Model, tea.Cmd) {
	model := &Model[Config]{
		ctx:   ctx,
		state: state,
		list:  app.LanguageFromMnemonicConfig(&app.MnemonicConfig{Language: language}),
		echo:  textinput.EchoPassword,

		style:    newStyle(tui.DefaultColorscheme()),
		bindings: newBindings(),
		help:     help.New(),

		next: next,
	}

	if ctx.IsSet(tui.MnemonicFlagName) {
		m := bip39.SplitMnemonic(strings.TrimSpace(ctx.String(tui.MnemonicFlagName)))
		return model.onMnemonic(mnemonic{m})
	}

	return model, model.initInputs()
}

func newInput() textinput.Model {
	input := textinput.New()
	input.EchoMode = textinput.EchoPassword
	input.Prompt = ""
	input.Width = 15
	input.KeyMap = inputBinding(textinput.DefaultKeyMap)
	return input
}

func (m *Model[Config]) initInputs() tea.Cmd {
	m.inputs = make([]textinput.Model, 1, mnemonicMaxWords)
	m.inputs[0] = newInput()
	m.focused = 0
	return m.inputs[0].Focus()
}

func (m *Model[Config]) Init() tea.Cmd {
	return nil
}

func (m *Model[Config]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errorMsg:
		m.errorMsg = msg
	case mnemonic:
		return m.onMnemonic(msg)
	case next:
		return m.onNext(msg)
	case prev:
		return m.onPrev(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.bindings.accept):
			return m.onConfirm()
		case key.Matches(msg, m.bindings.toggle):
			m.echo ^= 1
			return m, nil
		case key.Matches(msg, m.bindings.next, m.bindings.space):
			return m, Next(m.focused)
		case key.Matches(msg, m.bindings.prev):
			return m, Prev(m.focused)
		case key.Matches(msg, m.bindings.backspace):
			if m.inputs[m.focused].Value() == "" {
				return m, Prev(m.focused)
			}
		case key.Matches(msg, m.bindings.quit):
			return m, tui.QuitWithError(errors.New("mnemoinc input canceled"))
		}
	}

	return m, m.updateInputs(msg)
}

func (m *Model[Config]) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	for i := 0; i < len(m.inputs); i++ {
		var updateCmd, pasteCmd tea.Cmd
		m.inputs[i], updateCmd = m.inputs[i].Update(msg)
		if updateCmd != nil {
			cmds = append(cmds, updateCmd)
		}
		pasteCmd = m.onPaste(i)
		if pasteCmd != nil {
			cmds = append(cmds, pasteCmd)
		}
	}

	m.checkMnemonic()
	return tea.Batch(cmds...)
}

func (m *Model[Config]) onPaste(index int) tea.Cmd {
	words := strings.Fields(strings.TrimSpace(m.inputs[index].Value()))
	if len(words) < 2 {
		return nil
	}

	inputs := make([]textinput.Model, len(m.inputs)+len(words)-1)
	for i := 0; i < index; i++ {
		inputs[i] = m.inputs[i]
		inputs[i].Blur()
	}

	idx := index
	for i := 0; i < len(words); i++ {
		input := newInput()
		input.SetValue(words[i])
		inputs[idx] = input
		idx++
	}

	for i := index + 1; i < len(m.inputs); i++ {
		inputs[idx] = m.inputs[i]
		inputs[idx].Blur()
		idx++
	}

	m.inputs = inputs
	return Next(index + len(words) - 1)
}

func (m *Model[Config]) onMnemonic(msg mnemonic) (tea.Model, tea.Cmd) {
	m.state.WithMnemonic(msg.mnemonic, m.list)
	return m.next(m.ctx, m.state)
}

func (m *Model[Config]) checkMnemonic() {
	var (
		words      = len(m.inputs)
		isInBounds = words <= mnemonicMaxWords && words >= mnemonicMinWords
		isMultiple = words%3 == 0
	)

	m.bindings.accept.SetEnabled(isInBounds && isMultiple)
}

func (m *Model[Config]) onConfirm() (tea.Model, tea.Cmd) {
	mnemonic := make([]string, len(m.inputs))
	for i := range m.inputs {
		mnemonic[i] = m.inputs[i].Value()
	}

	if err := bip39.ValidateMnemonic(mnemonic, m.list); err != nil {
		return m, Error(errors.New("invalid mnemonic"))
	}

	return m, Confirm(mnemonic)
}

func (m *Model[Config]) onNext(msg next) (tea.Model, tea.Cmd) {
	index := msg.prev
	if m.inputs[index].Value() == "" {
		return m, nil
	}

	m.inputs[index].Blur()
	if index+1 == len(m.inputs) {
		m.inputs = append(m.inputs, newInput())
	}

	m.focused = index + 1
	return m, m.inputs[index+1].Focus()
}

func (m *Model[Config]) onPrev(msg prev) (tea.Model, tea.Cmd) {
	index := msg.next
	if index == 0 {
		return m, nil
	}

	if m.inputs[index].Value() == "" {
		if index != len(m.inputs)-1 {
			return m, nil
		}
		m.inputs[index].Blur()
		m.inputs = m.inputs[:len(m.inputs)-1]
	} else {
		m.inputs[index].Blur()
	}

	m.focused = index - 1
	return m, m.inputs[index-1].Focus()
}

func (m *Model[Config]) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.style.title.Foreground(m.style.colors.Title).Render("Mnemonic"),
		m.style.container.Render(
			lipgloss.JoinHorizontal(lipgloss.Bottom, m.mnemonicColumnsView()...),
		),
		m.errorView(),
		m.help.View(m.bindings),
	)
}

func (m *Model[Config]) errorView() string {
	if m.errorMsg.err != nil {
		return m.style.error.Foreground(m.style.colors.Error).Render(
			fmt.Sprintf("[Error]: %s", m.errorMsg.err.Error()),
		)
	}
	return ""
}

func (m *Model[Config]) mnemonicColumnsView() []string {
	views := make([]string, mnemonicColumns)
	for column := range mnemonicColumns {
		views[column] = m.mnemonicColumnView(column)
	}

	return views
}

func (m *Model[Config]) mnemonicColumnView(column int) string {
	var (
		rows  = len(m.inputs)/mnemonicColumns + 1
		views = make([]string, rows)
	)

	for row := range rows {
		index := row*mnemonicColumns + column
		if index >= len(m.inputs) {
			break
		}

		views[row] = m.wordView(index)
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (m *Model[Config]) wordView(index int) string {
	input := m.inputs[index]
	input.TextStyle = m.style.word.Foreground(m.style.colors.Title)
	if input.Focused() {
		input.TextStyle = m.style.word.Foreground(m.style.colors.Accent)
	}

	input.EchoMode = m.echo
	return m.style.wordContainer.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.style.index.Foreground(m.style.colors.Text).Render(strconv.Itoa(index+1)),
			input.View(),
		),
	)
}
