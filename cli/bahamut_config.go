//go:build !ethereum

package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/types"
	keystore "github.com/viwet/GoKeystoreV4"
)

const AppName = "Bahamut chain Staking CLI"

// DepositConfig stores all deposit related params
type DepositConfig struct {
	StartIndex uint32 `json:"start_index"`
	Number     uint32 `json:"number"`

	ChainConfig *config.ChainConfig `json:"chain_config,omitempty"`

	Amounts             *IndexedConfigWithDefault[uint64]  `json:"amounts,omitempty"`
	ContractAddresses   *IndexedConfig[Address]            `json:"contract_addresses,omitempty"`
	WithdrawalAddresses *IndexedConfigWithDefault[Address] `json:"withdrawal_addresses,omitempty"`

	Directory string `json:"directory"`

	KeystoreKeyDerivationFunction string `json:"kdf,omitempty"`
}

func validateDepositConfig(cfg *DepositConfig) error {
	if cfg.Number == 0 {
		return fmt.Errorf("cannot generate zero deposits")
	}

	if len(cfg.ChainConfig.GenesisForkVersion) != config.ForkVersionLength {
		return fmt.Errorf("invalid fork version length")
	}

	var (
		from = cfg.StartIndex
		to   = cfg.StartIndex + cfg.Number
	)

	if err := validateAmounts(cfg.Amounts, from, to); err != nil {
		return err
	}

	if err := validateContractAddresses(cfg.ContractAddresses, from, to); err != nil {
		return err
	}

	if err := validateWithdrawalAddresses(cfg.WithdrawalAddresses, from, to); err != nil {
		return err
	}

	if cfg.KeystoreKeyDerivationFunction != "" {
		switch cfg.KeystoreKeyDerivationFunction {
		case keystore.ScryptName, keystore.PBKDF2Name:
		default:
			return fmt.Errorf("invalid key derivation function (only scrypt and pbkdf2 allowed)")
		}
	}

	return nil
}

func newDepositConfigFromFlags(ctx *cli.Context) (*DepositConfig, error) {
	var (
		startIndex          = uint32(ctx.Uint(StartIndexFlag.Name))
		number              = uint32(ctx.Uint(NumberFlag.Name))
		amounts             = ctx.StringSlice(AmountsFlag.Name)
		contractAddresses   = ctx.StringSlice(ContractAddressesFlag.Name)
		withdrawalAddresses = ctx.StringSlice(WithdrawalAddressesFlag.Name)
		directory           = ctx.String(DirectoryFlag.Name)
		keystoreKDF         = ctx.String(KeystoreKDFFlag.Name)

		from, to = startIndex, startIndex + number
	)

	chainConfig, err := parseChainConfig(ctx)
	if err != nil {
		return nil, err
	}

	amountsConfig, err := parseAmounts(amounts, from, to)
	if err != nil {
		return nil, err
	}

	contractAddressesConfig, err := parseContractAddresses(contractAddresses, from, to)
	if err != nil {
		return nil, err
	}

	withdrawalAddressesConfig, err := parseWithdrawalAddresses(withdrawalAddresses, from, to)
	if err != nil {
		return nil, err
	}

	return &DepositConfig{
		StartIndex:                    startIndex,
		Number:                        number,
		ChainConfig:                   chainConfig,
		Amounts:                       amountsConfig,
		ContractAddresses:             contractAddressesConfig,
		WithdrawalAddresses:           withdrawalAddressesConfig,
		Directory:                     directory,
		KeystoreKeyDerivationFunction: keystoreKDF,
	}, nil
}

func parseContractAddresses(contracts []string, from, to uint32) (*IndexedConfig[Address], error) {
	config := &IndexedConfig[Address]{
		Config: make(map[uint32]Address),
	}

	unique := make(map[[20]byte]uint32)
	for _, contract := range contracts {
		values := strings.Split(contract, ":")
		if len(values) != 2 {
			if to-from == 1 {
				contract, err := ParseAddress(values[0])
				if err != nil {
					return nil, err
				}

				config.Config[from] = contract
				continue
			}
			return nil, fmt.Errorf("cannot process `contracts` flag value")
		}

		index, err := ParseIndex(values[0], from, to)
		if err != nil {
			return nil, err
		}

		contract, err := ParseAddress(values[1])
		if err != nil {
			return nil, err
		}

		if i, ok := unique[[20]byte(contract)]; ok && i != uint32(index) {
			return nil, fmt.Errorf("contract addresses must be unique")
		}

		config.Config[uint32(index)] = contract
	}

	return config, nil
}

// DepositOptions from config for given key index
func (cfg *DepositConfig) DepositOptions(index uint32) types.DepositOptions {
	var options []types.DepositOption
	if cfg.Amounts != nil {
		if amount := cfg.Amounts.Get(index); amount != 0 {
			options = append(options, types.WithAmount(amount))
		}
	}
	if cfg.WithdrawalAddresses != nil {
		if withdrawalAddress := cfg.WithdrawalAddresses.Get(index); len(withdrawalAddress) != 0 {
			options = append(options, types.WithWithdrawalAddress(withdrawalAddress))
		}
	}
	if cfg.ContractAddresses != nil {
		if contractAddress, ok := cfg.ContractAddresses.Get(index); ok && len(contractAddress) != 0 {
			options = append(options, types.WithContract(contractAddress))
		}
	}

	return options
}

func validateContractAddresses(contracts *IndexedConfig[Address], from, to uint32) error {
	if contracts == nil {
		return nil
	}

	unique := make(map[[20]byte]uint32)
	for index, contract := range contracts.Config {
		if !IsValidIndex(index, from, to) {
			return fmt.Errorf(
				"invalid contracts index: expected index between %d and %d, but got %d",
				from,
				to,
				index,
			)
		}

		if !IsValidAddress(contract) {
			return fmt.Errorf(
				"invalid contract address: expected length %d, but got %d",
				config.ExecutionAddressLength,
				len(contract),
			)
		}

		if i, ok := unique[[20]byte(contract)]; ok && i != index {
			return fmt.Errorf("contract addresses must be unique")
		}
	}

	return nil
}
