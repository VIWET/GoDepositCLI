package cli

import "github.com/urfave/cli/v2"

var (
	StartIndexFlag = &cli.UintFlag{
		Name:     "start-index",
		Usage:    "Key index from which deposits generation starts",
		Value:    0,
		Aliases:  []string{"index", "start", "from", "i"},
	}

	NumberFlag = &cli.UintFlag{
		Name:     "number",
		Usage:    "Number of deposit to generate",
		Value:    1,
		Aliases:  []string{"num", "n"},
	}

	AmountsFlag = &cli.StringSliceFlag{
		Name:     "amounts",
		Usage:    "Amounts to deposit",
		Aliases:  []string{"amount", "a"},
	}

	WithdrawalAddressesFlag = &cli.StringSliceFlag{
		Name:     "withdrawal-addresses",
		Usage:    "Withdrawal addresses to deposit with",
		Aliases:  []string{"withdraw-to", "wc"},
	}

	DirectoryFlag = &cli.StringFlag{
		Name:     "directory",
		Usage:    "Directory to store generated deposit data and keystores",
		Value:    "keys",
		Aliases:  []string{"dir", "d"},
	}

	KeystoreKDFFlag = &cli.StringFlag{
		Name:     "keystore-kdf",
		Usage:    "Key derivation function (scrypt, pbkdf2)",
	}

	DepositConfigFlag = &cli.StringFlag{
		Name:     "config",
		Usage:    "Path to deposit config",
		Aliases:  []string{"cfg"},
	}

	ChainNameFlag = &cli.StringFlag{
		Name:     "chain",
		Usage:    "Chain to deposit",
		Value:    "mainnet",
		Aliases:  []string{"network"},
	}

	ChainGenesisForkVersion = &cli.StringFlag{
		Name:     "genesis-fork",
		Usage:    "Chain genesis fork version",
	}

	ChainGenesisValidatorsRoot = &cli.StringFlag{
		Name:     "genesis-validators-root",
		Usage:    "Chain genesis validators root",
	}

	PasswordFlag = &cli.StringFlag{
		Name:     "password",
		Usage:    "Keystore password",
	}

	MnemonicFlag = &cli.StringFlag{
		Name:     "mnemonic",
		Usage:    "Seed phrase",
	}

	MnemonicLanguageFlag = &cli.StringFlag{
		Name:     "mnemonic-language",
		Usage:    "Language of seed phrase",
		Value:    "english",
		Aliases:  []string{"language", "lang", "l"},
	}

	MnemonicBitlenFlag = &cli.UintFlag{
		Name:     "mnemonic-bitlen",
		Usage:    "Strength of seed generated",
		Value:    256,
		Aliases:  []string{"strength", "bitlen", "bl", "s"},
	}

	MnemonicConfigFlag = &cli.StringFlag{
		Name:     "mnemonic-config",
		Usage:    "Path to mnemonic config",
		Aliases:  []string{"mnemonic-cfg", "mcfg"},
	}

	BLSToExecutionConfigFlag = &cli.StringFlag{
		Name:     "bls-to-execution-config",
		Usage:    "Path to bls to execution config",
		Aliases:  []string{"blscfg"},
	}

	ValidatorIndicesFlag = &cli.StringSliceFlag{
		Name:     "validator-indices",
		Usage:    "Indices to generate bls to execution messages",
		Aliases:  []string{"indices", "vi"},
	}
)
