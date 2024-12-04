package app

import (
	"bytes"
	"fmt"

	"github.com/viwet/GoDepositCLI/config"
)

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
