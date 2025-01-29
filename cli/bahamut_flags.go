//go:build !ethereum

package cli

import "github.com/urfave/cli/v2"

var ContractAddressesFlag = &cli.StringSliceFlag{
	Name:    "contracts",
	Usage:   "Contract address to deposit with",
	Aliases: []string{"contract", "c"},
}

var (
	depositNewMnemonicFlags = []cli.Flag{
		ConfigFlag,

		StartIndexFlag,
		NumberFlag,
		ChainNameFlag,
		ChainGenesisForkVersionFlag,
		ChainGenesisValidatorsRootFlag,
		MnemonicLanguageFlag,
		MnemonicBitlenFlag,
		DirectoryFlag,
		AmountsFlag,
		WithdrawalAddressesFlag,
		ContractAddressesFlag,
		KeystoreKDFFlag,

		PasswordFlag,
	}

	depositExistingMnemonicFlags = []cli.Flag{
		ConfigFlag,

		StartIndexFlag,
		NumberFlag,
		ChainNameFlag,
		ChainGenesisForkVersionFlag,
		ChainGenesisValidatorsRootFlag,
		MnemonicFlag,
		MnemonicLanguageFlag,
		DirectoryFlag,
		AmountsFlag,
		WithdrawalAddressesFlag,
		ContractAddressesFlag,
		KeystoreKDFFlag,

		PasswordFlag,
	}
)
