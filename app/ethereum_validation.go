//go:build ethereum

package app

import (
	"fmt"

	keystore "github.com/viwet/GoKeystoreV4"
)

func ensureDepositConfigIsValid(cfg *DepositConfig) error {
	if cfg.Config == nil {
		cfg.Config = new(Config)
	}

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
