package cli

import (
	"runtime"

	"github.com/urfave/cli/v2"
)

var (
	NonInteractiveFlag = &cli.BoolFlag{
		Name:  "non-interactive",
		Value: false,
	}

	EngineWorkersFlag = &cli.IntFlag{
		Name:  "engine-workers",
		Value: runtime.NumCPU(),
	}

	StartIndexFlag = &cli.UintFlag{
		Name:    "start-index",
		Usage:   "Key index from which deposits generation starts",
		Aliases: []string{"index", "start", "from", "i"},
	}

	NumberFlag = &cli.UintFlag{
		Name:    "number",
		Usage:   "Number of deposit to generate",
		Aliases: []string{"num", "n"},
	}

	ChainNameFlag = &cli.StringFlag{
		Name:    "chain",
		Usage:   "Chain to deposit",
		Aliases: []string{"network"},
	}

	ChainGenesisForkVersionFlag = &cli.StringFlag{
		Name:  "genesis-fork",
		Usage: "Chain genesis fork version",
	}

	ChainGenesisValidatorsRootFlag = &cli.StringFlag{
		Name:  "genesis-validators-root",
		Usage: "Chain genesis validators root",
	}

	MnemonicFlag = &cli.StringFlag{
		Name:  "mnemonic",
		Usage: "Seed phrase",
	}

	PasswordFlag = &cli.StringFlag{
		Name:  "password",
		Usage: "Keystore password",
	}

	MnemonicLanguageFlag = &cli.StringFlag{
		Name:    "language",
		Usage:   "Language of seed phrase",
		Aliases: []string{"lang", "l"},
	}

	MnemonicBitlenFlag = &cli.UintFlag{
		Name:    "bitlen",
		Usage:   "Strength of seed generated",
		Aliases: []string{"strength"},
	}

	DirectoryFlag = &cli.StringFlag{
		Name:    "directory",
		Usage:   "Directory to store generated deposit data and keystores",
		Aliases: []string{"dir", "d"},
	}

	AmountsFlag = &cli.StringSliceFlag{
		Name:    "amounts",
		Usage:   "Amounts to deposit",
		Aliases: []string{"amount", "a"},
	}

	WithdrawalAddressesFlag = &cli.StringSliceFlag{
		Name:    "withdrawal-addresses",
		Usage:   "Withdrawal addresses to deposit with",
		Aliases: []string{"withdrawal-address", "withdraw-to", "w"},
	}

	KeystoreKDFFlag = &cli.StringFlag{
		Name:  "kdf",
		Usage: "Key derivation function (scrypt, pbkdf2)",
	}

	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Usage:   "Path to config file",
		Aliases: []string{"cfg"},
	}

	ValidatorIndicesFlag = &cli.StringSliceFlag{
		Name:    "validator-indices",
		Usage:   "Indices to generate bls to execution messages",
		Aliases: []string{"validator-index", "indices", "vi"},
	}
)

var blsToExecutionFlags = []cli.Flag{
	ConfigFlag,

	StartIndexFlag,
	NumberFlag,
	ChainNameFlag,
	ChainGenesisForkVersionFlag,
	ChainGenesisValidatorsRootFlag,
	MnemonicLanguageFlag,
	DirectoryFlag,
	WithdrawalAddressesFlag,
	ValidatorIndicesFlag,
}
