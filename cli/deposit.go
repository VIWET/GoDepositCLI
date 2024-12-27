package cli

import (
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
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

	return generateDeposits(state.WithPassword(password))
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
		app.NewState(cfg).
			WithMnemonic(mnemonic, app.LanguageFromMnemonicConfig(cfg.MnemonicConfig)).
			WithPassword(password),
	)
}

func generateDeposits(state *app.State[app.DepositConfig]) error {
	deposits, keystores, err := app.GenerateDeposits(state)
	if err != nil {
		return err
	}

	cfg := state.Config()
	if err := ensureDirectoryExist(cfg.Directory); err != nil {
		return err
	}

	if err := saveDeposits(deposits, cfg.Directory); err != nil {
		return err
	}

	if err := saveKeystores(keystores, cfg.Directory); err != nil {
		return err
	}

	return nil
}
