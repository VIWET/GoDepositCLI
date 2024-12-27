package cli

import (
	"github.com/urfave/cli/v2"
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

	state := app.NewState(cfg).
		WithMnemonic(mnemonic, app.LanguageFromMnemonicConfig(cfg.MnemonicConfig))

	return generateBLSToExecution(state)
}

func generateBLSToExecution(state *app.State[app.BLSToExecutionConfig]) error {
	messages, err := app.GenerateBLSToExecutionMessages(state)
	if err != nil {
		return err
	}

	cfg := state.Config()
	if err := ensureDirectoryExist(cfg.Directory); err != nil {
		return err
	}

	if err := saveBLSToExecution(messages, cfg.Directory); err != nil {
		return err
	}

	return nil
}
