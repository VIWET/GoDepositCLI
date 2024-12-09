package cli

import "github.com/urfave/cli/v2"

var (
	StartIndexFlag = &cli.UintFlag{
		Name:    "start-index",
		Usage:   "Key index from which deposits generation starts",
		Value:   0,
		Aliases: []string{"index", "start", "from", "i"},
	}

	NumberFlag = &cli.UintFlag{
		Name:    "number",
		Usage:   "Number of deposit to generate",
		Value:   1,
		Aliases: []string{"num", "n"},
	}

	ChainNameFlag = &cli.StringFlag{
		Name:    "chain",
		Usage:   "Chain to deposit",
		Value:   "mainnet",
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
		Value:   "english",
		Aliases: []string{"lang", "l"},
	}

	MnemonicBitlenFlag = &cli.UintFlag{
		Name:    "bitlen",
		Usage:   "Strength of seed generated",
		Value:   256,
		Aliases: []string{"strength"},
	}

	DirectoryFlag = &cli.StringFlag{
		Name:    "directory",
		Usage:   "Directory to store generated deposit data and keystores",
		Value:   "keys",
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
)
