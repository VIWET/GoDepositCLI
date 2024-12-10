//go:build !ethereum

package app

import "fmt"

// EnsureDepositConfigIsValid validates all deposit generation related configurations
func EnsureDepositConfigIsValid(cfg *DepositConfig) error {
	if cfg.Config == nil {
		cfg.Config = new(Config)
	}

	if err := ensureConfigIsValid(cfg.Config, true); err != nil {
		return err
	}

	var (
		from = cfg.StartIndex
		to   = cfg.StartIndex + cfg.Number
	)

	if err := ensureAmountsConfigIsValid(cfg.Amounts, from, to); err != nil {
		return err
	}

	if err := ensureWithdrawalAddressesConfigIsValid(cfg.WithdrawalAddresses, from, to); err != nil {
		return err
	}

	if err := ensureContractAddressesConfigIsValid(cfg.ContractAddresses, from, to); err != nil {
		return err
	}

	if err := ensureKeyDerivationFunctionIsValid(cfg); err != nil {
		return err
	}

	return nil
}

func ensureContractAddressesConfigIsValid(cfg *IndexedConfig[Address], from, to uint32) error {
	if cfg == nil {
		return nil
	}

	unique := make(map[Address]uint32)
	for index, contract := range cfg.Config {
		if !IsValidIndex(index, from, to) {
			return fmt.Errorf(
				"invalid contract addresses config: key index should be between %d and %d, but got %d",
				from,
				to,
				index,
			)
		}

		if existingIndex, ok := unique[contract]; ok && existingIndex != index {
			return fmt.Errorf(
				"invalid contract addresses config: %d and %d have the same contract address %s",
				index,
				existingIndex,
				contract.ToChecksumHex(),
			)
		}

		unique[contract] = index
	}

	return nil
}
