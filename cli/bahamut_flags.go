//go:build !ethereum

package cli

import "github.com/urfave/cli/v2"

var ContractAddressesFlag = &cli.StringSliceFlag{
	Name:    "contracts",
	Usage:   "Contract address to deposit with",
	Aliases: []string{"contract", "c"},
}

var depositFlags = []cli.Flag{
	DepositConfigFlag,

	StartIndexFlag,
	NumberFlag,
	AmountsFlag,
	WithdrawalAddressesFlag,
	ContractAddressesFlag,
	DirectoryFlag,
	KeystoreKDFFlag,
	ChainNameFlag,
	ChainGenesisForkVersion,
	ChainGenesisValidatorsRoot,
	PasswordFlag,
}
