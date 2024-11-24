//go:build !ethereum

package cli

import "github.com/urfave/cli/v2"

var ContractAddressesFlag = &cli.StringSliceFlag{
	Name:     "contracts",
	Category: "Deposit",
	Usage:    "Contract address to deposit with",
	Aliases:  []string{"contract", "c"},
}
