package cli

import (
	"encoding/json"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/helpers"
)

func NewDepositConfigFromCLI(ctx *cli.Context) (*app.DepositConfig, error) {
	if ctx.IsSet(ConfigFlag.Name) {
		return newDepositConfigFromFile(ctx)
	}

	return newDepositConfigFromFlags(ctx)
}

func newDepositConfigFromFile(ctx *cli.Context) (*app.DepositConfig, error) {
	path := ctx.String(ConfigFlag.Name)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := new(app.DepositConfig)
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	if err := app.EnsureDepositConfigIsValid(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

type DepositConfigBuilder struct {
	amounts             []string
	withdrawalAddresses []string
	contractAddresses   []string

	cfg *app.DepositConfig
}

func NewDepositConfigBuilder() *DepositConfigBuilder {
	return &DepositConfigBuilder{
		cfg: new(app.DepositConfig),
	}
}

func (b *DepositConfigBuilder) StartIndex(index uint32) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(app.Config)
	}
	b.cfg.StartIndex = index
	return b
}

func (b *DepositConfigBuilder) Number(number uint32) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(app.Config)
	}
	b.cfg.Number = number
	return b
}

func (b *DepositConfigBuilder) Chain(name string) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(app.Config)
	}
	if b.cfg.ChainConfig == nil {
		b.cfg.ChainConfig = new(config.ChainConfig)
	}
	b.cfg.ChainConfig.Name = name
	return b
}

func (b *DepositConfigBuilder) GenesisForkVersion(forkVersion []byte) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(app.Config)
	}
	if b.cfg.ChainConfig == nil {
		b.cfg.ChainConfig = new(config.ChainConfig)
	}
	b.cfg.ChainConfig.GenesisForkVersion = forkVersion
	return b
}

func (b *DepositConfigBuilder) GenesisValidatorsRoot(validatorsRoot []byte) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(app.Config)
	}
	if b.cfg.ChainConfig == nil {
		b.cfg.ChainConfig = new(config.ChainConfig)
	}
	b.cfg.ChainConfig.GenesisValidatorsRoot = validatorsRoot
	return b
}

func (b *DepositConfigBuilder) MnemonicLanguage(language string) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(app.Config)
	}
	if b.cfg.MnemonicConfig == nil {
		b.cfg.MnemonicConfig = new(app.MnemonicConfig)
	}
	b.cfg.MnemonicConfig.Language = language
	return b
}

func (b *DepositConfigBuilder) MnemonicBitlen(bitlen uint) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(app.Config)
	}
	if b.cfg.MnemonicConfig == nil {
		b.cfg.MnemonicConfig = new(app.MnemonicConfig)
	}
	b.cfg.MnemonicConfig.Bitlen = bitlen
	return b
}

func (b *DepositConfigBuilder) Directory(directory string) *DepositConfigBuilder {
	if b.cfg.Config == nil {
		b.cfg.Config = new(app.Config)
	}
	b.cfg.Directory = directory
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

func (b *DepositConfigBuilder) Build() (*app.DepositConfig, error) {
	if err := b.build(); err != nil {
		return nil, err
	}

	if err := app.EnsureDepositConfigIsValid(b.cfg); err != nil {
		return nil, err
	}

	return b.cfg, nil
}

func (b *DepositConfigBuilder) buildAmounts() error {
	if len(b.amounts) == 0 {
		return nil
	}

	b.cfg.Amounts = &app.IndexedConfigWithDefault[app.Amount]{
		IndexedConfig: app.IndexedConfig[app.Amount]{
			Config: make(map[uint32]app.Amount),
		},
	}

	onDefault := func(amount string) error {
		var a app.Amount
		if err := a.FromString(amount); err != nil {
			return err
		}

		b.cfg.Amounts.Default = a
		return nil
	}

	onIndexed := func(index uint32, amount string) error {
		var a app.Amount
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

	b.cfg.WithdrawalAddresses = &app.IndexedConfigWithDefault[app.Address]{
		IndexedConfig: app.IndexedConfig[app.Address]{
			Config: make(map[uint32]app.Address),
		},
	}

	onDefault := func(address string) error {
		var a app.Address
		if err := a.FromHex(address); err != nil {
			return err
		}

		b.cfg.WithdrawalAddresses.Default = a
		return nil
	}

	onIndexed := func(index uint32, address string) error {
		var a app.Address
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
