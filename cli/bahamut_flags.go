//go:build !ethereum

package cli

import "github.com/urfave/cli/v3"

const NetworkName = "Bahamut chain"

var ContractAddressesFlag = &cli.StringSliceFlag{
	Name:     "contracts",
	Category: "Validator options",
	Usage:    "Specify contract addresses for deposits",
	Aliases:  []string{"contract", "c"},
}

var (
	depositNewMnemonicFlags = []cli.Flag{
		// Mnemonic option
		MnemonicBitlenFlag,
		// Validator options
		AmountsFlag,
		WithdrawalAddressesFlag,
		ContractAddressesFlag,
		// Keystore options
		KeystoreKDFFlag,
		PasswordFlag,
	}

	depositExistingMnemonicFlags = []cli.Flag{
		// Mnemonic input
		MnemonicFlag,
		// Validator options
		AmountsFlag,
		WithdrawalAddressesFlag,
		ContractAddressesFlag,
		// Keystore options
		PasswordFlag,
	}
)
