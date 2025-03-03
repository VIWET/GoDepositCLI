package cli

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/io"
	"github.com/viwet/GoDepositCLI/tui"
	"github.com/viwet/GoDepositCLI/tui/components/mnemonic"
	mnemonicInput "github.com/viwet/GoDepositCLI/tui/components/mnemonic_input"
)

func GenerateDepositsFromNewMnemonic(ctx context.Context, cmd *cli.Command) error {
	cfg, err := NewDepositConfigFromCLI(cmd)
	if err != nil {
		return err
	}

	state := app.NewState(cfg)
	if cmd.Bool(NonInteractiveFlag.Name) {
		return generateDepositsFromNewMnemonicNonInteractive(ctx, cmd, state)
	}

	return tui.Run(ctx, cmd, state, mnemonic.New)
}

func generateDepositsFromNewMnemonicNonInteractive(
	ctx context.Context,
	cmd *cli.Command,
	state *app.State[app.DepositConfig],
) error {
	mnemonic, list, err := app.GenerateMnemonic(state)
	if err != nil {
		return err
	}

	state.WithMnemonic(mnemonic, list)

	password, err := ReadPassword(cmd)
	if err != nil {
		return err
	}

	ShowMnemonic(state)

	return generateDeposits(ctx, state.WithPassword(password))
}

func GenerateDepositsFromExistingMnemonic(ctx context.Context, cmd *cli.Command) error {
	cfg, err := NewDepositConfigFromCLI(cmd)
	if err != nil {
		return err
	}

	state := app.NewState(cfg)
	if cmd.Bool(NonInteractiveFlag.Name) {
		return generateDepositsFromExistingMnemonicNonInteractive(ctx, cmd, state)
	}

	return tui.Run(ctx, cmd, state, mnemonicInput.NewDepositMnemonicInput())
}

func generateDepositsFromExistingMnemonicNonInteractive(
	ctx context.Context,
	cmd *cli.Command,
	state *app.State[app.DepositConfig],
) error {
	mnemonic, err := ReadMnemonic(cmd)
	if err != nil {
		return err
	}

	password, err := ReadPassword(cmd)
	if err != nil {
		return err
	}

	return generateDeposits(
		ctx,
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
