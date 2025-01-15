package cli

import (
	"context"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/io"
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

	return generateBLSToExecution(ctx.Context, state)
}

func generateBLSToExecution(ctx context.Context, state *app.State[app.BLSToExecutionConfig]) error {
	messages, err := app.NewBLSToExecutionEngine(state).Generate(ctx)
	if err != nil {
		return err
	}

	cfg := state.Config()
	if err := io.EnsureDirectoryExist(cfg.Directory); err != nil {
		return err
	}

	if err := io.SaveBLSToExecution(messages, cfg.Directory); err != nil {
		return err
	}

	return nil
}
