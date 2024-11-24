//go:build !ethereum

package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/types"
)

const AppName = "Bahamut chain Staking CLI"

// DepositConfig stores all deposit related params
type DepositConfig struct {
	StartIndex uint32 `json:"start_index"`
	Number     uint32 `json:"number"`

	ChainConfig *config.ChainConfig `json:"chain_config,omitempty"`

	Amounts             *IndexedConfigWithDefault[uint64]  `json:"amount,omitempty"`
	ContractAddresses   *IndexedConfig[Address]            `json:"contract_addresses,omitempty"`
	WithdrawalAddresses *IndexedConfigWithDefault[Address] `json:"withdrawal_addresses,omitempty"`

	Directory string `json:"directory"`

	KeystoreKeyDerivationFunction string `json:"kdf,omitempty"`
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

		config.Config[uint32(index)] = contract
	}

	return config, nil
}

// DepositOptions from config for given key index
func (cfg *DepositConfig) DepositOptions(index uint32) types.DepositOptions {
	var options []types.DepositOption
	if amount := cfg.Amounts.Get(index); amount != 0 {
		options = append(options, types.WithAmount(amount))
	}
	if withdrawalAddress := cfg.WithdrawalAddresses.Get(index); len(withdrawalAddress) != 0 {
		options = append(options, types.WithWithdrawalAddress(withdrawalAddress))
	}
	if contractAddress, ok := cfg.ContractAddresses.Get(index); ok && len(contractAddress) != 0 {
		options = append(options, types.WithContract(contractAddress))
	}

	return options
}

var depositFlags = []cli.Flag{
	DepositConfigFlag,

	StartIndexFlag,
	NumberFlag,
	AmountsFlag,
	WithdrawalAddressesFlag,
	ContractAddressesFlag,
	DirectoryFlag,
	KeystoreKDFFlag,
	ChainNameFlag,
	ChainGenesisForkVersion,
	ChainGenesisValidatorsRoot,
	PasswordFlag,
}
