package mnemonic

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui/v2"
)

const mnemonicColumns = 3

type Model struct {
	ctx   *cli.Context
	state *app.State[app.DepositConfig]

	bindings bindings
	style    style

	show bool
	help help.Model
}

func New(ctx *cli.Context, state *app.State[app.DepositConfig]) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if mnemonic := state.Mnemonic(); len(mnemonic) == 0 {
		cmd = generateMnemonic(state)
	}

	return Model{
		ctx:      ctx,
		state:    state,
		bindings: newBindings(),
		style:    newStyle(tui.DefaultColorscheme()),
		help:     help.New(),
	}, cmd
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *Mnemonic:
		if msg.err != nil {
			return m, tui.QuitWithError(msg.err)
		}
		m.state.WithMnemonic(msg.mnemonic, msg.list)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.bindings.toggle):
			m.show = !m.show
		case key.Matches(msg, m.bindings.accept):
			return m, tui.Quit()
		case key.Matches(msg, m.bindings.language):
			return NewLanguage(m.ctx, m.state)
		case key.Matches(msg, m.bindings.quit):
			return m, tui.QuitWithError(errors.New("mnemonic wasn't accepted"))
		}
	}

	return m, nil
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.style.title.Foreground(m.style.colors.White).Render("Mnemonic"),
		m.style.container.Render(
			lipgloss.JoinHorizontal(lipgloss.Bottom, m.mnemonicColumnsView()...),
		),
		m.help.View(m.bindings),
	)
}

func (m Model) mnemonicColumnsView() []string {
	views := make([]string, mnemonicColumns)
	for column := range mnemonicColumns {
		views[column] = m.mnemonicColumnView(column)
	}

	return views
}

func (m Model) mnemonicColumnView(column int) string {
	var (
		mnemonic = m.state.Mnemonic()

		rows  = len(mnemonic) / mnemonicColumns
		views = make([]string, rows)
	)

	for row := range rows {
		var (
			index = row*mnemonicColumns + column
			word  = mnemonic[index]
		)

		if !m.show {
			word = strings.Repeat("*", utf8.RuneCountInString(word))
		}

		views[row] = m.wordView(word, index+1)
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (m Model) wordView(word string, index int) string {
	return m.style.wordContainer.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.style.index.Foreground(m.style.colors.Black).Render(strconv.Itoa(index)),
			m.style.word.Foreground(m.style.colors.White).Render(word),
		),
	)
}
