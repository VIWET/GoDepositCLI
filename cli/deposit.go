package cli

import (
	"context"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/io"
	"github.com/viwet/GoDepositCLI/tui"
	"github.com/viwet/GoDepositCLI/tui/components/mnemonic"
	mnemonicInput "github.com/viwet/GoDepositCLI/tui/components/mnemonic_input"
)

func GenerateDepositsFromNewMnemonic(ctx *cli.Context) error {
	cfg, err := NewDepositConfigFromCLI(ctx)
	if err != nil {
		return err
	}

	state := app.NewState(cfg)
	if ctx.Bool(NonInteractiveFlag.Name) {
		return generateDepositsFromNewMnemonicNonInteractive(ctx, state)
	}

	return tui.Run(ctx, state, mnemonic.New)
}

func generateDepositsFromNewMnemonicNonInteractive(ctx *cli.Context, state *app.State[app.DepositConfig]) error {
	mnemonic, list, err := app.GenerateMnemonic(state)
	if err != nil {
		return err
	}

	state.WithMnemonic(mnemonic, list)

	password, err := ReadPassword(ctx)
	if err != nil {
		return err
	}

	ShowMnemonic(state)

	return generateDeposits(ctx.Context, state.WithPassword(password))
}

func GenerateDepositsFromExistingMnemonic(ctx *cli.Context) error {
	cfg, err := NewDepositConfigFromCLI(ctx)
	if err != nil {
		return err
	}

	state := app.NewState(cfg)
	if ctx.Bool(NonInteractiveFlag.Name) {
		return generateDepositsFromExistingMnemonicNonInteractive(ctx, state)
	}

	return tui.Run(ctx, state, mnemonicInput.NewDepositMnemonicInput())
}

func generateDepositsFromExistingMnemonicNonInteractive(ctx *cli.Context, state *app.State[app.DepositConfig]) error {
	mnemonic, err := ReadMnemonic(ctx)
	if err != nil {
		return err
	}

	password, err := ReadPassword(ctx)
	if err != nil {
		return err
	}

	return generateDeposits(
		ctx.Context,
		state.
			WithMnemonic(mnemonic, app.LanguageFromMnemonicConfig(state.Config().MnemonicConfig)).
			WithPassword(password),
	)
}

func generateDeposits(ctx context.Context, state *app.State[app.DepositConfig]) error {
	deposits, keystores, err := app.NewDepositEngine(state).Generate(ctx)
	if err != nil {
		return err
	}

	cfg := state.Config()
	if err := io.EnsureDirectoryExist(cfg.Directory); err != nil {
		return err
	}

	if err := io.SaveDeposits(deposits, cfg.Directory); err != nil {
		return err
	}

	if err := io.SaveKeystores(keystores, cfg.Directory); err != nil {
		return err
	}

	return nil
}
