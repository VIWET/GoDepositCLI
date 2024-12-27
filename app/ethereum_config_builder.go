//go:build ethereum

package app

func (b *DepositConfigBuilder) build() error {
	if err := b.buildAmounts(); err != nil {
		return err
	}

	if err := b.buildWithdrawalAddresses(); err != nil {
		return err
	}

	return nil
}
