//go:build ethereum

package app

// EnsureDepositConfigIsValid validates all deposit generation related configurations
func EnsureDepositConfigIsValid(cfg *DepositConfig) error {
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

	if err := ensureKeyDerivationFunctionIsValid(cfg); err != nil {
		return err
	}

	return nil
}
