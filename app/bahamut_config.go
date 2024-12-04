//go:build !ethereum

package app

import (
	"fmt"

	keystore "github.com/viwet/GoKeystoreV4"
)

// DepositConfig stores all deposit generation related data
type DepositConfig struct {
	*Config

	Amounts             *IndexedConfigWithDefault[uint64]  `json:"amounts,omitempty"`
	WithdrawalAddresses *IndexedConfigWithDefault[Address] `json:"withdrawal_addresses,omitempty"`

	ContractAddresses *IndexedConfig[Address] `json:"contract_addresses,omitempty"`

	KeystoreKeyDerivationFunction string `json:"kdf,omitempty"`
}

func ensureDepositConfigIsValid(cfg *DepositConfig) error {
	if err := ensureConfigIsValid(cfg.Config); err != nil {
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

	if cfg.KeystoreKeyDerivationFunction == "" {
		cfg.KeystoreKeyDerivationFunction = keystore.ScryptName
	} else {
		switch cfg.KeystoreKeyDerivationFunction {
		case keystore.ScryptName, keystore.PBKDF2Name:
		default:
			return fmt.Errorf("invalid deposit config: %w", ErrInvalidKDF)
		}
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
