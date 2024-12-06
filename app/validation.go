package app

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/viwet/GoDepositCLI/config"
	keystore "github.com/viwet/GoKeystoreV4"
)

// EnsureBLSToExecutionConfigIsValid validates all bls to execution generation related configurations
func EnsureBLSToExecutionConfigIsValid(cfg *BLSToExecutionConfig) error {
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

	if err := ensureWithdrawalAddressesConfigIsValidBLS(cfg.WithdrawalAddresses, from, to); err != nil {
		return err
	}

	if err := ensureValidatorIndicesConfigIsValid(cfg.ValidatorIndices, from, to); err != nil {
		return err
	}

	return nil
}

func ensureConfigIsValid(cfg *Config) error {
	if cfg.Number == 0 {
		cfg.Number = 1
	}

	if cfg.ChainConfig == nil {
		cfg.ChainConfig = config.MainnetConfig()
	}
	if err := ensureChainConfigIsValid(cfg.ChainConfig); err != nil {
		return fmt.Errorf("invalid chain config: %w", err)
	}

	if cfg.MnemonicConfig == nil {
		cfg.MnemonicConfig = new(MnemonicConfig)
	}
	if err := ensureMnemonicConfigIsValid(cfg.MnemonicConfig); err != nil {
		return fmt.Errorf("invalid mnemonic config: %w", err)
	}

	if cfg.Directory == "" {
		cfg.Directory = DefaultOutputDirectory
	}

	return nil
}

// TODO(viwet): make GenesisValidatorsRoot validation optional
func ensureChainConfigIsValid(cfg *config.ChainConfig) error {
	if knownConfig, ok := config.ConfigByNetworkName(cfg.Name); ok {
		return ensureKnownChainConfigIsValid(cfg, knownConfig)
	}

	if len(cfg.GenesisForkVersion) != config.ForkVersionLength {
		return ErrInvalidGenesisForkVersion
	}

	if len(cfg.GenesisValidatorsRoot) != config.HashLength {
		return ErrInvalidGenesisValidatorsRoot
	}

	return nil
}

func ensureKnownChainConfigIsValid(cfg, knownConfig *config.ChainConfig) error {
	if len(cfg.GenesisForkVersion) == 0 {
		cfg.GenesisForkVersion = knownConfig.GenesisForkVersion
	} else if !bytes.Equal(cfg.GenesisForkVersion, knownConfig.GenesisForkVersion) {
		return fmt.Errorf(
			"different genesis fork version on %s - want: 0x%x, got: 0x%x",
			cfg.Name,
			knownConfig.GenesisForkVersion,
			cfg.GenesisForkVersion,
		)
	}

	if len(cfg.GenesisValidatorsRoot) == 0 {
		cfg.GenesisValidatorsRoot = knownConfig.GenesisValidatorsRoot
	} else if !bytes.Equal(cfg.GenesisValidatorsRoot, knownConfig.GenesisValidatorsRoot) {
		return fmt.Errorf(
			"different genesis validators root on %s - want: 0x%x, got: 0x%x",
			cfg.Name,
			knownConfig.GenesisValidatorsRoot,
			cfg.GenesisValidatorsRoot,
		)
	}

	return nil
}

func ensureMnemonicConfigIsValid(cfg *MnemonicConfig) error {
	if cfg.Language == "" {
		cfg.Language = "english"
	} else if err := validateMnemonicLanguage(cfg.Language); err != nil {
		return err
	}

	if cfg.Bitlen == 0 {
		cfg.Bitlen = 256
	} else if err := validateMnemonicBitlen(cfg.Bitlen); err != nil {
		return err
	}

	return nil
}

func ensureAmountsConfigIsValid(cfg *IndexedConfigWithDefault[uint64], from, to uint32) error {
	if cfg == nil {
		return nil
	}

	if cfg.Default != 0 {
		if !IsValidAmount(cfg.Default) {
			return fmt.Errorf(
				"invalid default amount %d: %w",
				cfg.Default,
				ErrInvalidAmount,
			)
		}
	}

	for index, amount := range cfg.Config {
		if !IsValidIndex(index, from, to) {
			return fmt.Errorf(
				"invalid amount config: key index should be between %d and %d, but got %d",
				from,
				to,
				index,
			)
		}

		if !IsValidAmount(amount) {
			return fmt.Errorf(
				"invalid amount config: invalid amount at index %d (%d): %w",
				index,
				amount,
				ErrInvalidAmount,
			)
		}
	}

	return nil
}

func ensureWithdrawalAddressesConfigIsValid(cfg *IndexedConfigWithDefault[Address], from, to uint32) error {
	if cfg == nil {
		return nil
	}

	for index := range cfg.Config {
		if !IsValidIndex(index, from, to) {
			return fmt.Errorf(
				"invalid withdrawal addresses config: key index should be between %d and %d, but got %d",
				from,
				to,
				index,
			)
		}
	}

	return nil
}

func ensureKeyDerivationFunctionIsValid(cfg *DepositConfig) error {
	if cfg.KeystoreKeyDerivationFunction != "" {
		switch strings.ToLower(cfg.KeystoreKeyDerivationFunction) {
		case keystore.ScryptName, keystore.PBKDF2Name:
		default:
			return fmt.Errorf("invalid deposit config: %w", ErrInvalidKDF)
		}
	}

	return nil
}

func ensureWithdrawalAddressesConfigIsValidBLS(cfg *IndexedConfigWithDefault[Address], from, to uint32) error {
	if cfg == nil {
		return ErrNoWithdrawalAddresses
	}

	hasDefault := cfg.Default != zeroAddress
	for index := range cfg.Config {
		if !IsValidIndex(index, from, to) {
			return fmt.Errorf(
				"invalid withdrawal addresses config: key index should be between %d and %d, but got %d",
				from,
				to,
				index,
			)
		}
	}

	defaultCount := (to - from) - uint32(len(cfg.Config))
	if !hasDefault && defaultCount > 0 {
		if defaultCount == to-from {
			return ErrNoWithdrawalAddresses
		}

		missed := make([]uint32, 0, defaultCount)
		for index := from; index < to; index++ {
			if _, ok := cfg.Config[index]; !ok {
				missed = append(missed, index)
			}
		}

		return fmt.Errorf("no withdrawal addresses for key indices: %v", missed)
	}

	return nil
}

func ensureValidatorIndicesConfigIsValid(cfg *IndexedConfig[uint64], from, to uint32) error {
	if cfg == nil {
		return ErrNoValidatorIndices
	}

	unique := make(map[uint64]uint32)
	for index, validatorIndex := range cfg.Config {
		if !IsValidIndex(index, from, to) {
			return fmt.Errorf(
				"invalid validator indices config: key index should be between %d and %d, but got %d",
				from,
				to,
				index,
			)
		}

		if existingIndex, ok := unique[validatorIndex]; ok && existingIndex != index {
			return fmt.Errorf(
				"invalid validator indices config: %d and %d have the same validator index %d",
				index,
				existingIndex,
				validatorIndex,
			)
		}

		unique[validatorIndex] = index
	}

	if uint32(len(unique)) < to-from {
		missedCount := to - from - uint32(len(unique))
		missed := make([]uint32, 0, missedCount)
		for index := from; index < to; index++ {
			if _, ok := cfg.Config[index]; !ok {
				missed = append(missed, index)
			}
		}

		return fmt.Errorf("no validator indices for key indices: %v", missed)
	}

	return nil
}

// IsValidAmount returns false if amount less than MinDepositAmount or greater than MaxDepositAmount
func IsValidAmount(amount uint64) bool {
	return config.MinDepositAmount <= amount &&
		amount <= config.MaxDepositAmount &&
		amount%uint64(config.GweiPerEther) == 0
}

// IsValidIndex returns false if index less than from or greater than to
func IsValidIndex(index, from, to uint32) bool {
	return from <= index && index < to
}
