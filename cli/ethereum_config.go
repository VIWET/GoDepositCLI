//go:build ethereum

package cli

import (
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/types"
)

// DepositConfig stores all deposit related params
type DepositConfig struct {
	StartIndex uint32 `json:"start_index"`
	Number     uint32 `json:"number"`

	ChainConfig *config.ChainConfig `json:"chain_config,omitempty"`

	Amounts             *IndexedConfigWithDefault[uint64]  `json:"amount,omitempty"`
	WithdrawalAddresses *IndexedConfigWithDefault[Address] `json:"withdrawal_addresses,omitempty"`

	Directory string `json:"directory"`

	KeystoreKeyDerivationFunction string `json:"kdf,omitempty"`
}

func newDepositConfigFromFlags(ctx *cli.Context) (*DepositConfig, error) {
	var (
		startIndex          = uint32(ctx.Uint(StartIndexFlag.Name))
		number              = uint32(ctx.Uint(NumberFlag.Name))
		amounts             = ctx.StringSlice(AmountsFlag.Name)
		withdrawalAddresses = ctx.StringSlice(WithdrawalAddressesFlag.Name)
		directory           = ctx.String(DirectoryFlag.Name)
		keystoreKDF         = ctx.String(KeystoreKDFFlag.Name)

		from, to = startIndex, startIndex + number
	)

	chainConfig, err := parseChainConfig(ctx)
	if err != nil {
		return nil, err
	}

	amountsConfig, err := parseAmounts(amounts, from, to)
	if err != nil {
		return nil, err
	}

	withdrawalAddressesConfig, err := parseWithdrawalAddresses(withdrawalAddresses, from, to)
	if err != nil {
		return nil, err
	}

	return &DepositConfig{
		StartIndex:                    startIndex,
		Number:                        number,
		ChainConfig:                   chainConfig,
		Amounts:                       amountsConfig,
		WithdrawalAddresses:           withdrawalAddressesConfig,
		Directory:                     directory,
		KeystoreKeyDerivationFunction: keystoreKDF,
	}, nil
}

// DepositOptions from config for given key index
func (cfg *DepositConfig) DepositOptions(index uint32) types.DepositOptions {
	var options []types.DepositOption
	if amount := cfg.Amounts.Get(index); amount != 0 {
		options = append(options, types.WithAmount(amount))
	}
	if withdrawalAddress := cfg.WithdrawalAddresses.Get(index); len(withdrawalAddress) != 0 {
		options = append(options, types.WithWithdrawalAddress(withdrawalAddress))
	}

	return options
}
