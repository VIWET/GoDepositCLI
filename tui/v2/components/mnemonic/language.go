package mnemonic

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui/v2/components/menu"
)

var languages = [...]string{
	"English",
	"Chinese Simplified",
	"Chinese Traditional",
	"Czech",
	"French",
	"Italian",
	"Japanese",
	"Korean",
	"Portuguese",
	"Spanish",
}

func NewLanguage(ctx *cli.Context, state *app.State[app.DepositConfig]) (tea.Model, tea.Cmd) {
	return menu.New("Language", generateLanguageOptions(ctx, state)...), nil
}

func generateLanguageOptions(ctx *cli.Context, state *app.State[app.DepositConfig]) []menu.Option {
	options := make([]menu.Option, len(languages))
	for i, lang := range languages {
		options[i] = menu.NewOption(lang, func() (tea.Model, tea.Cmd) {
			state.WithMnemonic(nil, nil)
			state.Config().MnemonicConfig.Language = lang
			return New(ctx, state)
		})
	}

	return options
}
