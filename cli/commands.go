package cli

import "github.com/urfave/cli/v2"

var (
	NewMnemonicCommand = &cli.Command{
		Name:    "new-mnemonic",
		Aliases: []string{"new"},
		Usage:   "Generate new mnemonic and deposits",
		Flags:   depositNewMnemonicFlags,
		Action:  GenerateDepositsFromNewMnemonic,
	}

	ExistingMnemonicCommand = &cli.Command{
		Name:    "existing-mnemonic",
		Aliases: []string{"existing"},
		Usage:   "Generate deposits using existing mnemonic",
		Flags:   depositExistingMnemonicFlags,
		Action:  GenerateDepositsFromExistingMnemonic,
	}

	DepositCommand = &cli.Command{
		Name:  "deposit",
		Usage: "Generate deposits using new mnemonic or existing one",
		Subcommands: []*cli.Command{
			NewMnemonicCommand,
			ExistingMnemonicCommand,
		},
	}
)
