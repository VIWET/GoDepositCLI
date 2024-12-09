//go:build !ethereum

package cli

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
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

	b.cfg.ContractAddresses = &app.IndexedConfig[app.Address]{
		Config: make(map[uint32]app.Address),
	}

	onDefault := func(address string) error {
		return fmt.Errorf("invalid contract addresses config: default contract address %s is not allowed", address)
	}

	onIndexed := func(index uint32, address string) error {
		var a app.Address
		if err := a.FromHex(address); err != nil {
			return err
		}

		b.cfg.ContractAddresses.Config[index] = a
		return nil
	}

	if err := helpers.ParseIndexedValues(onDefault, onIndexed, b.amounts...); err != nil {
		return err
	}

	return nil
}

func newDepositConfigFromFlags(ctx *cli.Context) (*app.DepositConfig, error) {
	builder := NewDepositConfigBuilder()

	builder.StartIndex(uint32(ctx.Uint(StartIndexFlag.Name)))
	builder.Number(uint32(ctx.Uint(NumberFlag.Name)))

	if ctx.IsSet(ChainNameFlag.Name) {
		builder.Chain(ctx.String(ChainNameFlag.Name))
	}

	if ctx.IsSet(ChainGenesisForkVersionFlag.Name) {
		forkVersion, err := hex.DecodeString(strings.TrimPrefix(ctx.String(ChainGenesisForkVersionFlag.Name), "0x"))
		if err != nil {
			return nil, err
		}
		builder.GenesisForkVersion(forkVersion)
	}

	if ctx.IsSet(ChainGenesisForkVersionFlag.Name) {
		forkVersion, err := hex.DecodeString(
			strings.TrimPrefix(
				ctx.String(ChainGenesisForkVersionFlag.Name),
				"0x",
			),
		)
		if err != nil {
			return nil, err
		}
		builder.GenesisForkVersion(forkVersion)
	}

	if ctx.IsSet(ChainGenesisValidatorsRootFlag.Name) {
		forkVersion, err := hex.DecodeString(
			strings.TrimPrefix(
				ctx.String(ChainGenesisValidatorsRootFlag.Name),
				"0x",
			),
		)
		if err != nil {
			return nil, err
		}
		builder.GenesisForkVersion(forkVersion)
	}

	builder.MnemonicLanguage(ctx.String(MnemonicLanguageFlag.Name))
	builder.MnemonicBitlen(ctx.Uint(MnemonicBitlenFlag.Name))
	builder.Directory(ctx.String(DirectoryFlag.Name))
	builder.Amounts(ctx.StringSlice(AmountsFlag.Name)...)
	builder.WithdrawalAddresses(ctx.StringSlice(WithdrawalAddressesFlag.Name)...)
	builder.ContractAddresses(ctx.StringSlice(ContractAddressesFlag.Name)...)
	builder.KeystoreKDF(ctx.String(KeystoreKDFFlag.Name))

	return builder.Build()
}
