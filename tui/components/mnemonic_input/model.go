package mnemonic_input

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui"
	"github.com/viwet/GoDepositCLI/tui/components/password"
)

type Model[Config app.ConfigConstraint] struct {
	ctx   *cli.Context
	state *app.State[Config]

	language string
	next     tui.NewModel[Config]
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
	if ctx.IsSet(tui.MnemonicFlagName) {
		mnemonic := bip39.SplitMnemonic(strings.TrimSpace(ctx.String(tui.MnemonicFlagName)))
		state.WithMnemonic(mnemonic, app.LanguageFromMnemonicConfig(&app.MnemonicConfig{Language: language}))
		return next(ctx, state)
	}

	model := &Model[Config]{
		ctx:   ctx,
		state: state,

		language: language,

		next: next,
	}

	return model, tui.QuitWithError(errors.New("mnemonic input model is unimplemented - coming soon"))
}

func (m *Model[Config]) Init() tea.Cmd {
	return nil
}

func (m *Model[Config]) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m *Model[Config]) View() string {
	return ""
}
