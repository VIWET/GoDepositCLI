package cli

import (
	"context"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/io"
)

func GenerateDepositsFromNewMnemonic(ctx *cli.Context) error {
	cfg, err := NewDepositConfigFromCLI(ctx)
	if err != nil {
		return err
	}

	state := app.NewState(cfg)
	mnemonic, list, err := app.GenerateMnemonic(state)
	if err != nil {
		return err
	}

	state.WithMnemonic(mnemonic, list)
	if err := ShowMnemonic(ctx, state); err != nil {
		return err
	}

	password, err := ReadPassword(ctx)
	if err != nil {
		return err
	}

	return generateDeposits(ctx.Context, state.WithPassword(password))
}

func GenerateDepositsFromExistingMnemonic(ctx *cli.Context) error {
	cfg, err := NewDepositConfigFromCLI(ctx)
	if err != nil {
		return err
	}

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
		app.NewState(cfg).
			WithMnemonic(mnemonic, app.LanguageFromMnemonicConfig(cfg.MnemonicConfig)).
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
