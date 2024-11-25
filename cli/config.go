package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/config"

	keystore "github.com/viwet/GoKeystoreV4"
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
	if len(validatorsRoot) != config.HashLength && ctx.Command.Name == "bls-to-execution" {
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

// Crypto returns KeystoreV4 crypto modules from config
func (cfg *DepositConfig) CryptoOptions() (keystore.CryptoOptions, error) {
	var options []keystore.CryptoOption
	if kdf := cfg.KeystoreKeyDerivationFunction; kdf != "" {
		switch strings.ToLower(kdf) {
		case keystore.PBKDF2Name:
			options = append(options, keystore.WithKDF(keystore.NewPBKDF2()))
		case keystore.ScryptName:
			options = append(options, keystore.WithKDF(keystore.NewScrypt()))
		default:
			return nil, fmt.Errorf("unknown kdf function: %s", kdf)
		}
	}

	return options, nil
}

// MnemonicConfig config
type MnemonicConfig struct {
	Language string `json:"language"`
	Bitlen   uint   `json:"bitlen"`
}

// NewMnemonicConfigFromCLI return config from file if config file provided or from flags
func NewMnemonicConfigFromCLI(ctx *cli.Context) (*MnemonicConfig, error) {
	if filepath := ctx.String(MnemonicConfigFlag.Name); filepath != "" {
		return newMnemonicConifgFromFile(filepath)
	}

	return newMnemonicConfigFromFlags(ctx)
}

func newMnemonicConifgFromFile(filepath string) (*MnemonicConfig, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := new(MnemonicConfig)
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func newMnemonicConfigFromFlags(ctx *cli.Context) (*MnemonicConfig, error) {
	var (
		language = ctx.String(MnemonicLanguageFlag.Name)
		bitlen   = ctx.Uint(MnemonicBitlenFlag.Name)
	)

	return &MnemonicConfig{
		Language: language,
		Bitlen:   bitlen,
	}, nil
}

// BLSToExecution config
type BLSToExecutionConfig struct {
	StartIndex uint32 `json:"start_index"`
	Number     uint32 `json:"number"`

	ChainConfig *config.ChainConfig `json:"chain_config,omitempty"`

	ValidatorIndices    *IndexedConfig[uint64]             `json:"validator_indices,omitempty"`
	WithdrawalAddresses *IndexedConfigWithDefault[Address] `json:"withdrawal_addresses,omitempty"`

	Directory string `json:"directory"`
}

// NewBLSToExecutionConfigFromCLI return config from file if config file provided or from flags
func NewBLSToExecutionConfigFromCLI(ctx *cli.Context) (*BLSToExecutionConfig, error) {
	if filepath := ctx.String(BLSToExecutionConfigFlag.Name); filepath != "" {
		return newBLSToExecutionConfigFromFile(filepath)
	}

	return newBLSToExecutionConfigFromFlags(ctx)
}

func newBLSToExecutionConfigFromFile(filepath string) (*BLSToExecutionConfig, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := new(BLSToExecutionConfig)
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	if cfg.ChainConfig == nil {
		cfg.ChainConfig = config.MainnetConfig()
	}

	return cfg, nil
}

func newBLSToExecutionConfigFromFlags(ctx *cli.Context) (*BLSToExecutionConfig, error) {
	var (
		startIndex          = uint32(ctx.Uint(StartIndexFlag.Name))
		number              = uint32(ctx.Uint(NumberFlag.Name))
		indices             = ctx.StringSlice(ValidatorIndicesFlag.Name)
		withdrawalAddresses = ctx.StringSlice(WithdrawalAddressesFlag.Name)
		directory           = ctx.String(DirectoryFlag.Name)

		from, to = startIndex, startIndex + number
	)

	chainConfig, err := parseChainConfig(ctx)
	if err != nil {
		return nil, err
	}

	validatorIndicesConfig, err := parseValidatorIndices(indices, from, to)
	if err != nil {
		return nil, err
	}

	withdrawalAddressesConfig, err := parseWithdrawalAddresses(withdrawalAddresses, from, to)
	if err != nil {
		return nil, err
	}

	if len(withdrawalAddressesConfig.Default) != config.ExecutionAddressLength && len(withdrawalAddressesConfig.Config) != int(number) {
		return nil, fmt.Errorf("withdrawal credentials must be set for all provided validators")
	}

	return &BLSToExecutionConfig{
		StartIndex:          startIndex,
		Number:              number,
		ChainConfig:         chainConfig,
		ValidatorIndices:    validatorIndicesConfig,
		WithdrawalAddresses: withdrawalAddressesConfig,
		Directory:           directory,
	}, nil
}

func parseValidatorIndices(indices []string, from, to uint32) (*IndexedConfig[uint64], error) {
	config := &IndexedConfig[uint64]{
		Config: make(map[uint32]uint64),
	}

	unique := make(map[uint64]uint32)
	for _, index := range indices {
		values := strings.Split(index, ":")
		if len(values) != 2 {
			if to-from == 1 {
				validator, err := ParseValidatorIndex(values[0])
				if err != nil {
					return nil, err
				}

				config.Config[from] = validator
				continue
			}

			return nil, fmt.Errorf("cannot process `validator-indices` flag value")
		}

		index, err := ParseIndex(values[0], from, to)
		if err != nil {
			return nil, err
		}

		validator, err := ParseValidatorIndex(values[1])
		if err != nil {
			return nil, err
		}

		if i, ok := unique[validator]; ok && i != uint32(index) {
			return nil, fmt.Errorf("validator indices must be unique")
		}

		config.Config[uint32(index)] = validator
	}

	return config, nil
}
