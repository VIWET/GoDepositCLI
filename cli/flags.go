package cli

import (
	"runtime"

	"github.com/urfave/cli/v3"
)

var (
	NonInteractiveFlag = &cli.BoolFlag{
		Name:     "non-interactive",
		Category: "App options",
		Usage:    "Run in non-interactive mode",
		Value:    false,
	}

	EngineWorkersFlag = &cli.IntFlag{
		Name:     "engine-workers",
		Category: "App options",
		Usage:    "Number of workers for processing deposits",
		Value:    int64(runtime.NumCPU()),
	}

	StartIndexFlag = &cli.UintFlag{
		Name:     "start-index",
		Category: "Generation options",
		Usage:    "Starting key index for deposit generation",
		Aliases:  []string{"index", "start", "from", "i"},
	}

	NumberFlag = &cli.UintFlag{
		Name:     "number",
		Category: "Generation options",
		Usage:    "Number of deposits to generate",
		Aliases:  []string{"num", "n"},
	}

	ChainNameFlag = &cli.StringFlag{
		Name:     "chain",
		Category: "Chain options",
		Usage:    "Target blockchain network for the deposit",
		Aliases:  []string{"network"},
	}

	ChainGenesisForkVersionFlag = &cli.StringFlag{
		Name:     "genesis-fork",
		Category: "Chain options",
		Usage:    "Genesis fork version of the blockchain",
	}

	ChainGenesisValidatorsRootFlag = &cli.StringFlag{
		Name:     "genesis-validators-root",
		Category: "Chain options",
		Usage:    "Root of the genesis validators set",
	}

	MnemonicFlag = &cli.StringFlag{
		Name:     "mnemonic",
		Category: "Mnemonic options",
		Usage:    "Seed phrase for key generation",
	}

	PasswordFlag = &cli.StringFlag{
		Name:     "password",
		Category: "Keystore options",
		Usage:    "Password for encrypting the keystore",
	}

	MnemonicLanguageFlag = &cli.StringFlag{
		Name:     "language",
		Category: "Mnemonic options",
		Usage:    "Language of the seed phrase",
		Aliases:  []string{"lang", "l"},
	}

	MnemonicBitlenFlag = &cli.UintFlag{
		Name:     "bitlen",
		Category: "Mnemonic options",
		Usage:    "Bit length (strength) of the generated seed",
		Aliases:  []string{"strength"},
	}

	DirectoryFlag = &cli.StringFlag{
		Name:     "directory",
		Category: "Generation options",
		Usage:    "Directory to store deposit data and keystores",
		Aliases:  []string{"dir", "d"},
	}

	AmountsFlag = &cli.StringSliceFlag{
		Name:     "amounts",
		Category: "Validator options",
		Usage:    "Deposit amounts",
		Aliases:  []string{"amount", "a"},
	}

	WithdrawalAddressesFlag = &cli.StringSliceFlag{
		Name:     "withdrawal-addresses",
		Category: "Validator options",
		Usage:    "Withdrawal addresses for the deposits",
		Aliases:  []string{"withdraw-to", "w"},
	}

	KeystoreKDFFlag = &cli.StringFlag{
		Name:     "kdf",
		Category: "Keystore options",
		Usage:    "Key derivation function (scrypt or pbkdf2)",
	}

	ConfigFlag = &cli.StringFlag{
		Name:     "config",
		Category: "Generation options",
		Usage:    "Path to the configuration file",
		Aliases:  []string{"cfg"},
	}

	ValidatorIndicesFlag = &cli.StringSliceFlag{
		Name:     "validator-indices",
		Category: "Validator options",
		Usage:    "Validator indices for generating BLS-to-execution messages",
		Aliases:  []string{"validator-index", "indices", "vi"},
	}
)

var sharedFlags = []cli.Flag{
	// App options
	NonInteractiveFlag,
	EngineWorkersFlag,
	// Generation options
	ConfigFlag,
	StartIndexFlag,
	NumberFlag,
	DirectoryFlag,
	// Chain options
	ChainNameFlag,
	ChainGenesisForkVersionFlag,
	ChainGenesisValidatorsRootFlag,
	// Mnemonic option
	MnemonicLanguageFlag,
}

var blsToExecutionFlags = []cli.Flag{
	// Mnemonic input
	MnemonicFlag,
	// Validator options
	WithdrawalAddressesFlag,
	ValidatorIndicesFlag,
}
