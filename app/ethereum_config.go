//go:build ethereum

package app

import "github.com/viwet/GoDepositCLI/types"

// DepositConfig stores all deposit generation related data
type DepositConfig struct {
	*Config

	Amounts             *IndexedConfigWithDefault[Amount]  `json:"amounts,omitempty"`
	WithdrawalAddresses *IndexedConfigWithDefault[Address] `json:"withdrawal_addresses,omitempty"`

	KeystoreKeyDerivationFunction string `json:"kdf,omitempty"`
}

func newDepositOptionsFromConfig(cfg *DepositConfig, index uint32) types.DepositOptions {
	var options []types.DepositOption
	if cfg.Amounts != nil {
		if amount := cfg.Amounts.Get(index); amount > 0 {
			options = append(options, types.WithAmount(amount.Gwei()))
		}
	}

	if cfg.WithdrawalAddresses != nil {
		if withdrawalAddress := cfg.WithdrawalAddresses.Get(index); withdrawalAddress != zeroAddress {
			options = append(options, types.WithWithdrawalAddress(withdrawalAddress[:]))
		}
	}

	return options
}
