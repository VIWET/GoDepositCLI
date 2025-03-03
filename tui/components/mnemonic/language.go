package mnemonic

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui/components/menu"
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

func NewLanguage(ctx context.Context, cmd *cli.Command, state *app.State[app.DepositConfig]) (tea.Model, tea.Cmd) {
	return menu.New("Language", generateLanguageOptions(ctx, cmd, state)...), nil
}

func generateLanguageOptions(ctx context.Context, cmd *cli.Command, state *app.State[app.DepositConfig]) []menu.Option {
	options := make([]menu.Option, len(languages))
	for i, lang := range languages {
		options[i] = menu.NewOption(lang, func() (tea.Model, tea.Cmd) {
			state.WithMnemonic(nil, nil)
			state.Config().MnemonicConfig.Language = lang
			return New(ctx, cmd, state)
		})
	}

	return options
}
