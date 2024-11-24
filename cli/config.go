package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/config"
)

// IndexedConfig stores values by key index
type IndexedConfig[T any] struct {
	Config map[uint32]T `json:"config"`
}

// Get value by key index
func (cfg *IndexedConfig[T]) Get(index uint32) (T, bool) {
	value, ok := cfg.Config[index]
	return value, ok
}

// IndexedConfigWithDefault stroes values by key index and default value
type IndexedConfigWithDefault[T any] struct {
	Default T `json:"default"`
	IndexedConfig[T]
}

// Get value by key index or default
func (cfg *IndexedConfigWithDefault[T]) Get(index uint32) T {
	value, ok := cfg.IndexedConfig.Get(index)
	if ok {
		return value
	}
	return cfg.Default
}

// NewDepositConfigFromCLI return config from file if config file provided or from flags
func NewDepositConfigFromCLI(ctx *cli.Context) (*DepositConfig, error) {
	if filepath := ctx.String(DepositConfigFlag.Name); filepath != "" {
		return newDepositConfigFromFile(filepath)
	}

	return newDepositConfigFromFlags(ctx)
}

func newDepositConfigFromFile(filepath string) (*DepositConfig, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := new(DepositConfig)
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	if cfg.ChainConfig == nil {
		cfg.ChainConfig = config.MainnetConfig()
	}

	return cfg, nil
}

func parseChainConfig(ctx *cli.Context) (*config.ChainConfig, error) {
	name := strings.ToLower(ctx.String(ChainNameFlag.Name))
	cfg, ok := config.ConfigByNetworkName(name)
	if ok {
		return cfg, nil
	}

	var (
		genesisForkVersion    = ctx.String(ChainGenesisForkVersion.Name)
		genesisValidatorsRoot = ctx.String(ChainGenesisValidatorsRoot.Name)
	)

	forkVersion, err := hex.DecodeString(strings.TrimPrefix(genesisForkVersion, "0x"))
	if err != nil {
		return nil, err
	}
	if len(forkVersion) != config.ForkVersionLength {
		return nil, fmt.Errorf("invalid fork version length")
	}

	validatorsRoot, err := hex.DecodeString(strings.TrimPrefix(genesisValidatorsRoot, "0x"))
	if err != nil {
		return nil, err
	}
	if len(validatorsRoot) != config.HashLength {
		return nil, fmt.Errorf("invalid genesis validators root length")
	}

	return &config.ChainConfig{
		Name:                  name,
		GenesisForkVersion:    forkVersion,
		GenesisValidatorsRoot: validatorsRoot,
	}, nil
}

func parseAmounts(amounts []string, from, to uint32) (*IndexedConfigWithDefault[uint64], error) {
	config := &IndexedConfigWithDefault[uint64]{
		IndexedConfig: IndexedConfig[uint64]{
			Config: make(map[uint32]uint64),
		},
	}

	for _, amount := range amounts {
		values := strings.Split(amount, ":")
		switch len(values) {
		case 1:
			amount, err := ParseGwei(values[0])
			if err != nil {
				return nil, err
			}

			config.Default = amount
		case 2:
			index, err := ParseIndex(values[0], from, to)
			if err != nil {
				return nil, err
			}

			amount, err := ParseGwei(values[1])
			if err != nil {
				return nil, err
			}

			config.Config[uint32(index)] = amount
		default:
			return nil, fmt.Errorf("cannot process `amount` flag value")
		}
	}

	return config, nil
}

func parseWithdrawalAddresses(addresses []string, from, to uint32) (*IndexedConfigWithDefault[Address], error) {
	config := &IndexedConfigWithDefault[Address]{
		IndexedConfig: IndexedConfig[Address]{
			Config: make(map[uint32]Address),
		},
	}

	for _, address := range addresses {
		values := strings.Split(address, ":")
		switch len(values) {
		case 1:
			address, err := ParseAddress(values[0])
			if err != nil {
				return nil, err
			}

			config.Default = address
		case 2:
			index, err := ParseIndex(values[0], from, to)
			if err != nil {
				return nil, err
			}

			address, err := ParseAddress(values[1])
			if err != nil {
				return nil, err
			}

			config.Config[uint32(index)] = address
		default:
			return nil, fmt.Errorf("cannot process `withdrawal-addresses` flag value")
		}

	}

	return config, nil
}
