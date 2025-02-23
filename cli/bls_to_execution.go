package cli

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/io"
	"github.com/viwet/GoDepositCLI/tui"
	mnemonicInput "github.com/viwet/GoDepositCLI/tui/components/mnemonic_input"
)

func GenerateBLSToExecution(ctx context.Context, cmd *cli.Command) error {
	cfg, err := NewBLSToExecutionConfigFromCLI(cmd)
	if err != nil {
		return err
	}

	state := app.NewState(cfg)
	if cmd.Bool(NonInteractiveFlag.Name) {
		return generateBLSToExecutionNonInteractive(ctx, cmd, state)
	}

	return tui.Run(ctx, cmd, state, mnemonicInput.NewBLSToExecutionMnemonicInput())
}

func generateBLSToExecutionNonInteractive(
	ctx context.Context,
	cmd *cli.Command,
	state *app.State[app.BLSToExecutionConfig],
) error {
	mnemonic, err := ReadMnemonic(cmd)
	if err != nil {
		return err
	}

	state.WithMnemonic(mnemonic, app.LanguageFromMnemonicConfig(state.Config().MnemonicConfig))
	return generateBLSToExecution(ctx, state)
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
