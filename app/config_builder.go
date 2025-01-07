package app

import (
	"fmt"
	"strconv"

	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/helpers"
)

type DepositConfigBuilder struct {
	amounts             []string
	withdrawalAddresses []string
	contractAddresses   []string

	cfg *DepositConfig
}

func NewDepositConfigBuilder() *DepositConfigBuilder {
	return &DepositConfigBuilder{
		cfg: new(DepositConfig),
	}
}

func (b *DepositConfigBuilder) StartIndex(index uint32) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	b.cfg.StartIndex = index
	return b
}

func (b *DepositConfigBuilder) Number(number uint32) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	b.cfg.Number = number
	return b
}

func (b *DepositConfigBuilder) Chain(name string) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	if b.cfg.ChainConfig == nil {
		b.cfg.ChainConfig = new(config.ChainConfig)
	}
	b.cfg.ChainConfig.Name = name
	return b
}

func (b *DepositConfigBuilder) GenesisForkVersion(forkVersion []byte) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	if b.cfg.ChainConfig == nil {
		b.cfg.ChainConfig = new(config.ChainConfig)
	}
	b.cfg.ChainConfig.GenesisForkVersion = forkVersion
	return b
}

func (b *DepositConfigBuilder) GenesisValidatorsRoot(validatorsRoot []byte) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	if b.cfg.ChainConfig == nil {
		b.cfg.ChainConfig = new(config.ChainConfig)
	}
	b.cfg.ChainConfig.GenesisValidatorsRoot = validatorsRoot
	return b
}

func (b *DepositConfigBuilder) MnemonicLanguage(language string) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	if b.cfg.MnemonicConfig == nil {
		b.cfg.MnemonicConfig = new(MnemonicConfig)
	}
	b.cfg.MnemonicConfig.Language = language
	return b
}

func (b *DepositConfigBuilder) MnemonicBitlen(bitlen uint) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	if b.cfg.MnemonicConfig == nil {
		b.cfg.MnemonicConfig = new(MnemonicConfig)
	}
	b.cfg.MnemonicConfig.Bitlen = bitlen
	return b
}

func (b *DepositConfigBuilder) Directory(directory string) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	b.cfg.Directory = directory
	return b
}

func (b *DepositConfigBuilder) EngineWorkers(workers int) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	b.cfg.EngineWorkers = max(workers, 1)
	return b
}

func (b *DepositConfigBuilder) Amounts(amounts ...string) *DepositConfigBuilder {
	b.amounts = append(b.amounts, amounts...)
	return b
}

func (b *DepositConfigBuilder) WithdrawalAddresses(addresses ...string) *DepositConfigBuilder {
	b.withdrawalAddresses = append(b.withdrawalAddresses, addresses...)
	return b
}

func (b *DepositConfigBuilder) KeystoreKDF(kdf string) *DepositConfigBuilder {
	b.cfg.KeystoreKeyDerivationFunction = kdf
	return b
}

func (b *DepositConfigBuilder) Build() (*DepositConfig, error) {
	if err := b.build(); err != nil {
		return nil, err
	}

	if err := EnsureDepositConfigIsValid(b.cfg); err != nil {
		return nil, err
	}

	return b.cfg, nil
}

type BLSToExecutionConfigBuilder struct {
	validatorIndices    []string
	withdrawalAddresses []string

	cfg *BLSToExecutionConfig
}

func NewBLSToExecutionConfigBuilder() *BLSToExecutionConfigBuilder {
	return &BLSToExecutionConfigBuilder{
		cfg: new(BLSToExecutionConfig),
	}
}

func (b *BLSToExecutionConfigBuilder) StartIndex(index uint32) *BLSToExecutionConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	b.cfg.StartIndex = index
	return b
}

func (b *BLSToExecutionConfigBuilder) Number(number uint32) *BLSToExecutionConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	b.cfg.Number = number
	return b
}

func (b *BLSToExecutionConfigBuilder) Chain(name string) *BLSToExecutionConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	if b.cfg.ChainConfig == nil {
		b.cfg.ChainConfig = new(config.ChainConfig)
	}
	b.cfg.ChainConfig.Name = name
	return b
}

func (b *BLSToExecutionConfigBuilder) GenesisForkVersion(forkVersion []byte) *BLSToExecutionConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	if b.cfg.ChainConfig == nil {
		b.cfg.ChainConfig = new(config.ChainConfig)
	}
	b.cfg.ChainConfig.GenesisForkVersion = forkVersion
	return b
}

func (b *BLSToExecutionConfigBuilder) GenesisValidatorsRoot(validatorsRoot []byte) *BLSToExecutionConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	if b.cfg.ChainConfig == nil {
		b.cfg.ChainConfig = new(config.ChainConfig)
	}
	b.cfg.ChainConfig.GenesisValidatorsRoot = validatorsRoot
	return b
}

func (b *BLSToExecutionConfigBuilder) MnemonicLanguage(language string) *BLSToExecutionConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	if b.cfg.MnemonicConfig == nil {
		b.cfg.MnemonicConfig = new(MnemonicConfig)
	}
	b.cfg.MnemonicConfig.Language = language
	return b
}

func (b *BLSToExecutionConfigBuilder) Directory(directory string) *BLSToExecutionConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	b.cfg.Directory = directory
	return b
}

func (b *BLSToExecutionConfigBuilder) EngineWorkers(workers int) *BLSToExecutionConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(Config)
	}
	b.cfg.EngineWorkers = max(workers, 1)
	return b
}

func (b *BLSToExecutionConfigBuilder) WithdrawalAddresses(addresses ...string) *BLSToExecutionConfigBuilder {
	b.withdrawalAddresses = append(b.withdrawalAddresses, addresses...)
	return b
}

func (b *BLSToExecutionConfigBuilder) ValidatorIndices(indices ...string) *BLSToExecutionConfigBuilder {
	b.validatorIndices = append(b.validatorIndices, indices...)
	return b
}

func (b *BLSToExecutionConfigBuilder) Build() (*BLSToExecutionConfig, error) {
	if err := b.build(); err != nil {
		return nil, err
	}

	if err := EnsureBLSToExecutionConfigIsValid(b.cfg); err != nil {
		return nil, err
	}

	return b.cfg, nil
}

func (b *DepositConfigBuilder) buildAmounts() error {
	if len(b.amounts) == 0 {
		return nil
	}

	b.cfg.Amounts = &IndexedConfigWithDefault[Amount]{
		IndexedConfig: IndexedConfig[Amount]{
			Config: make(map[uint32]Amount),
		},
	}

	onDefault := func(amount string) error {
		var a Amount
		if err := a.FromString(amount); err != nil {
			return err
		}

		b.cfg.Amounts.Default = a
		return nil
	}

	onIndexed := func(index uint32, amount string) error {
		var a Amount
		if err := a.FromString(amount); err != nil {
			return err
		}

		b.cfg.Amounts.Config[index] = a
		return nil
	}

	if err := helpers.ParseIndexedValues(onDefault, onIndexed, b.amounts...); err != nil {
		return err
	}

	return nil
}

func (b *DepositConfigBuilder) buildWithdrawalAddresses() error {
	if len(b.withdrawalAddresses) == 0 {
		return nil
	}

	b.cfg.WithdrawalAddresses = &IndexedConfigWithDefault[Address]{
		IndexedConfig: IndexedConfig[Address]{
			Config: make(map[uint32]Address),
		},
	}

	onDefault := func(address string) error {
		var a Address
		if err := a.FromHex(address); err != nil {
			return err
		}

		b.cfg.WithdrawalAddresses.Default = a
		return nil
	}

	onIndexed := func(index uint32, address string) error {
		var a Address
		if err := a.FromHex(address); err != nil {
			return err
		}

		b.cfg.WithdrawalAddresses.Config[index] = a
		return nil
	}

	if err := helpers.ParseIndexedValues(onDefault, onIndexed, b.withdrawalAddresses...); err != nil {
		return err
	}

	return nil
}

func (b *BLSToExecutionConfigBuilder) build() error {
	if err := b.buildWithdrawalAddresses(); err != nil {
		return err
	}

	if err := b.buildValidatorIndices(); err != nil {
		return err
	}

	return nil
}

func (b *BLSToExecutionConfigBuilder) buildWithdrawalAddresses() error {
	if len(b.withdrawalAddresses) == 0 {
		return nil
	}

	b.cfg.WithdrawalAddresses = &IndexedConfigWithDefault[Address]{
		IndexedConfig: IndexedConfig[Address]{
			Config: make(map[uint32]Address),
		},
	}

	onDefault := func(address string) error {
		var a Address
		if err := a.FromHex(address); err != nil {
			return err
		}

		b.cfg.WithdrawalAddresses.Default = a
		return nil
	}

	onIndexed := func(index uint32, address string) error {
		var a Address
		if err := a.FromHex(address); err != nil {
			return err
		}

		b.cfg.WithdrawalAddresses.Config[index] = a
		return nil
	}

	if err := helpers.ParseIndexedValues(onDefault, onIndexed, b.withdrawalAddresses...); err != nil {
		return err
	}

	return nil
}

func (b *BLSToExecutionConfigBuilder) buildValidatorIndices() error {
	if len(b.validatorIndices) == 0 {
		return nil
	}

	b.cfg.ValidatorIndices = &IndexedConfig[uint64]{
		Config: make(map[uint32]uint64),
	}

	onDefault := func(index string) error {
		return fmt.Errorf("invalid validator indices config: default validator index %s is not allowed", index)
	}

	onIndexed := func(index uint32, validatorIndex string) error {
		vIndex, err := strconv.ParseUint(validatorIndex, 10, 64)
		if err != nil {
			return err
		}

		b.cfg.ValidatorIndices.Config[index] = vIndex
		return nil
	}

	if err := helpers.ParseIndexedValues(onDefault, onIndexed, b.validatorIndices...); err != nil {
		return err
	}

	return nil
}
