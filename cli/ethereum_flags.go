//go:build ethereum

package cli

import "github.com/urfave/cli/v3"

const NetworkName = "Ethereum"

var (
	depositNewMnemonicFlags = []cli.Flag{
		// Mnemonic option
		MnemonicBitlenFlag,
		// Validator options
		AmountsFlag,
		WithdrawalAddressesFlag,
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
		// Keystore options
		KeystoreKDFFlag,
		PasswordFlag,
	}
)
