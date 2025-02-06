//go:build !ethereum

package app

import (
	"fmt"

	"github.com/viwet/GoDepositCLI/helpers"
)

func (b *DepositConfigBuilder) ContractAddresses(addresses ...string) *DepositConfigBuilder {
	b.contractAddresses = append(b.contractAddresses, addresses...)
	return b
}

func (b *DepositConfigBuilder) build() error {
	if err := b.buildAmounts(); err != nil {
		return err
	}

	if err := b.buildWithdrawalAddresses(); err != nil {
		return err
	}

	if err := b.buildContractAddresses(); err != nil {
		return err
	}

	return nil
}

func (b *DepositConfigBuilder) buildContractAddresses() error {
	if len(b.contractAddresses) == 0 {
		return nil
	}

	b.cfg.ContractAddresses = &IndexedConfig[Address]{
		Config: make(map[uint32]Address),
	}

	onDefault := func(address string) error {
		return fmt.Errorf("invalid contract addresses config: default contract address %s is not allowed", address)
	}

	onIndexed := func(index uint32, address string) error {
		var a Address
		if err := a.FromHex(address); err != nil {
			return err
		}

		b.cfg.ContractAddresses.Config[index] = a
		return nil
	}

	if err := helpers.ParseIndexedValues(onDefault, onIndexed, b.contractAddresses...); err != nil {
		return err
	}

	return nil
}
