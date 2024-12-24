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

	mnemonic, list, err := app.GenerateMnemonic(cfg.MnemonicConfig)
	if err != nil {
		return err
	}

	ShowMnemonic(mnemonic)
	password, err := ReadPassword(ctx)
	if err != nil {
		return err
	}

	state := app.NewState(cfg).
		WithMnemonic(mnemonic, list).
		WithPassword(password)

	return generateDeposits(state)
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

	state := app.NewState(cfg).
		WithMnemonic(mnemonic, app.LanguageFromMnemonicConfig(cfg.MnemonicConfig)).
		WithPassword(password)

	return generateDeposits(state)
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
