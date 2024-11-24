package cli

import "github.com/urfave/cli/v2"

var (
	StartIndexFlag = &cli.UintFlag{
		Name:     "start-index",
		Category: "Deposit",
		Usage:    "Key index from which deposits generation starts",
		Value:    0,
		Aliases:  []string{"index", "start", "from", "i"},
	}

	NumberFlag = &cli.UintFlag{
		Name:     "number",
		Category: "Deposit",
		Usage:    "Number of deposit to generate",
		Value:    1,
		Aliases:  []string{"num", "n"},
	}

	AmountsFlag = &cli.StringSliceFlag{
		Name:     "amounts",
		Category: "Deposit",
		Usage:    "Amounts to deposit",
		Aliases:  []string{"amount", "a"},
	}

	WithdrawalAddressesFlag = &cli.StringSliceFlag{
		Name:     "withdrawal-addresses",
		Category: "Deposit",
		Usage:    "Withdrawal addresses to deposit with",
		Aliases:  []string{"withdraw-to", "wc"},
	}

	DirectoryFlag = &cli.StringFlag{
		Name:     "directory",
		Category: "Deposit",
		Usage:    "Directory to store generated deposit data and keystores",
		Value:    "keys",
		Aliases:  []string{"dir", "d"},
	}

	KeystoreKDFFlag = &cli.StringFlag{
		Name:     "keystore-kdf",
		Category: "Deposit",
		Usage:    "Key derivation function (scrypt, pbkdf2)",
	}

	DepositConfigFlag = &cli.StringFlag{
		Name:     "config",
		Category: "Deposit",
		Usage:    "Path to deposit config",
		Aliases:  []string{"cfg"},
	}

	ChainNameFlag = &cli.StringFlag{
		Name:     "chain",
		Category: "Network",
		Usage:    "Chain to deposit",
		Value:    "mainnet",
		Aliases:  []string{"network"},
	}

	ChainGenesisForkVersion = &cli.StringFlag{
		Name:     "genesis-fork",
		Category: "Network",
		Usage:    "Chain genesis fork version",
	}

	ChainGenesisValidatorsRoot = &cli.StringFlag{
		Name:     "genesis-validators-root",
		Category: "Network",
		Usage:    "Chain genesis validators root",
	}

	PasswordFlag = &cli.StringFlag{
		Name:     "password",
		Category: "Deposit",
		Usage:    "Keystore password",
	}

	MnemonicFlag = &cli.StringFlag{
		Name:     "mnemonic",
		Category: "Mnemonic",
		Usage:    "Seed phrase",
	}

	MnemonicLanguageFlag = &cli.StringFlag{
		Name:     "mnemonic-language",
		Category: "Mnemonic",
		Usage:    "Language of seed phrase",
		Value:    "english",
		Aliases:  []string{"language", "lang", "l"},
	}

	MnemonicBitlenFlag = &cli.UintFlag{
		Name:     "mnemonic-bitlen",
		Category: "Mnemonic",
		Usage:    "Strength of seed generated",
		Value:    256,
		Aliases:  []string{"strength", "bitlen", "bl", "s"},
	}

	MnemonicConfigFlag = &cli.StringFlag{
		Name:     "mnemonic-config",
		Category: "Mnemonic",
		Usage:    "Path to mnemonic config",
		Aliases:  []string{"mnemonic-cfg", "mcfg"},
	}
)
