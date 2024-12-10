package cli

import (
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoBIP39/words"
	"github.com/viwet/GoDepositCLI/app"
)

func GenerateBLSToExecution(ctx *cli.Context) error {
	cfg, err := NewBLSToExecutionConfigFromCLI(ctx)
	if err != nil {
		return err
	}

	mnemonic, err := ReadMnemonic(ctx)
	if err != nil {
		return err
	}

	return generateBLSToExecution(
		cfg,
		mnemonic,
		app.LanguageFromMnemonicConfig(cfg.MnemonicConfig),
	)
}

func generateBLSToExecution(cfg *app.BLSToExecutionConfig, mnemonic []string, list words.List) error {
	messages, err := app.GenerateBLSToExecutionMessages(cfg, mnemonic, list)
	if err != nil {
		return err
	}

	if err := ensureDirectoryExist(cfg.Directory); err != nil {
		return err
	}

	if err := saveBLSToExecution(messages, cfg.Directory); err != nil {
		return err
	}

	return nil
}
