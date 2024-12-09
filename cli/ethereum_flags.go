//go:build ethereum

package cli

import "github.com/urfave/cli/v2"

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
		KeystoreKDFFlag,

		PasswordFlag,
	}
)
