package cli

import (
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoBIP39/words"
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

	return generateDeposits(
		cfg,
		mnemonic,
		list,
		password,
	)
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
		cfg,
		mnemonic,
		app.LanguageFromMnemonicConfig(cfg.MnemonicConfig),
		password,
	)
}

func generateDeposits(cfg *app.DepositConfig, mnemonic []string, list words.List, password string) error {
	deposits, keystores, err := app.GenerateDeposits(cfg, mnemonic, list, password)
	if err != nil {
		return err
	}

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
