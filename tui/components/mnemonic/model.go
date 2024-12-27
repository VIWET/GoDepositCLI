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
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui"
)

const columns = 3

type Model struct {
	state   *app.State[app.DepositConfig]
	show    bool
	binding bindings
	help    help.Model
}

func New(state *app.State[app.DepositConfig]) *Model {
	return &Model{
		state:   state,
		binding: newBindings(),
		help:    help.New(),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.binding.toggle):
			m.show = !m.show
		case key.Matches(msg, m.binding.accept):
			return m, tui.Quit()
		case key.Matches(msg, m.binding.quit):
			return m, tui.QuitWithError(errors.New("mnemonic wasn't accepted"))
		}
	}

	return m, nil
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		renderTitle("Mnemonic"),
		m.renderMnemonic(),
		renderHelp(m.help, m.binding),
	)
}

func renderTitle(title string) string {
	return titleStyle.Render(title)
}

func (m *Model) renderMnemonic() string {
	mnemonic, words := m.state.Mnemonic(), m.state.Words()
	return mnemonicSectionContainerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			renderLanguage(words().Language()),
			lipgloss.JoinHorizontal(
				lipgloss.Bottom,
				renderMnemonicColumn(mnemonic, 0, m.show),
				renderMnemonicColumn(mnemonic, 1, m.show),
				renderMnemonicColumn(mnemonic, 2, m.show),
			),
		),
	)
}

func renderLanguage(language string) string {
	return mnemonicLanguageSectionContainerStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			mnemonicLanguageStyle.Render("Language: "),
			mnemonicLanguageStyle.Italic(true).Render(language),
		),
	)
}

func renderMnemonicColumn(mnemonic []string, column int, show bool) string {
	var (
		rows  = len(mnemonic) / columns
		views = make([]string, rows)
	)

	for row := range rows {
		var (
			index = row*columns + column
			word  = mnemonic[index]
		)

		if !show {
			word = strings.Repeat("*", utf8.RuneCountInString(word))
		}

		views[row] = renderWordWithIndex(word, index+1)
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
			mnemonicWordStyle.Render(word),
		),
	)
}

func renderHelp(help help.Model, binding help.KeyMap) string {
	return help.View(binding)
}
