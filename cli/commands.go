package cli

import "github.com/urfave/cli/v2"

var (
	DepositNewMnemonicCommand = &cli.Command{
		Name:     "new-mnemonic",
		Aliases:  []string{"new"},
		Usage:    "Create deposits and new mnemonic",
		Category: "Deposit",
		Action:   GenerateDepositsNewMnemonic,
		Flags: []cli.Flag{
			MnemonicConfigFlag,

			MnemonicBitlenFlag,
			MnemonicLanguageFlag,
		},
	}

	DepositExistingMnemonicCommand = &cli.Command{
		Name:     "existing-mnemonic",
		Aliases:  []string{"existing"},
		Usage:    "Create deposits with existing mnemonic",
		Category: "Deposit",
		Action:   GenerateDepositsExistingMnemonic,
		Flags: []cli.Flag{
			MnemonicConfigFlag,

			MnemonicLanguageFlag,
			MnemonicFlag,
		},
	}

	DepositCommand = &cli.Command{
		Name:     "deposit",
		Usage:    "Create new deposits",
		Category: "Deposit",
		Subcommands: []*cli.Command{
			DepositNewMnemonicCommand,
			DepositExistingMnemonicCommand,
		},
		Flags: depositFlags,
	}

	BLSToExecutionCommand = &cli.Command{
		Name:     "bls-to-execution",
		Usage:    "Create BLSToExecution messages",
		Category: "BLS to Execution",
		Flags: []cli.Flag{
			MnemonicLanguageFlag,
			MnemonicFlag,

			BLSToExecutionConfigFlag,

			StartIndexFlag,
			NumberFlag,
			ValidatorIndicesFlag,
			WithdrawalAddressesFlag,
			DirectoryFlag,

			ChainNameFlag,
			ChainGenesisForkVersion,
			ChainGenesisValidatorsRoot,
		},
	}
)
