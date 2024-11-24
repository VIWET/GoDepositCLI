//go:build ethereum

package cli

import "github.com/urfave/cli/v2"

var depositFlags = []cli.Flag{
	DepositConfigFlag,

	StartIndexFlag,
	NumberFlag,
	AmountsFlag,
	WithdrawalAddressesFlag,
	DirectoryFlag,
	KeystoreKDFFlag,
	ChainNameFlag,
	ChainGenesisForkVersion,
	ChainGenesisValidatorsRoot,
	PasswordFlag,
}
