package app

import (
	"strings"

	"github.com/viwet/GoDepositCLI/config"
	keystore "github.com/viwet/GoKeystoreV4"
)

const DefaultOutputDirectory = "./validators_data"

// Config is a base config
type Config struct {
	StartIndex uint32 `json:"start_index"`
	Number     uint32 `json:"number"`

	ChainConfig    *config.ChainConfig `json:"chain_config,omitempty"`
	MnemonicConfig *MnemonicConfig     `json:"mnemonic_config,omitempty"`

	Directory string `json:"directory"`
}

// MnemonicConfig config
type MnemonicConfig struct {
	Language string `json:"language"`
	Bitlen   uint   `json:"bitlen"`
}

// IndexedConfig stores values by key index
type IndexedConfig[T any] struct {
	Config map[uint32]T `json:"config"`
}

// Get value by key index if exist
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

func newCryptoOptionsFromConfig(cfg *DepositConfig) keystore.CryptoOptions {
	var options []keystore.CryptoOption
	if kdf := cfg.KeystoreKeyDerivationFunction; kdf != "" {
		switch strings.ToLower(kdf) {
		case keystore.PBKDF2Name:
			options = append(options, keystore.WithKDF(keystore.NewPBKDF2()))
		case keystore.ScryptName:
			options = append(options, keystore.WithKDF(keystore.NewScrypt()))
		default:
			// Config is assumed validated
			panic(ErrInvalidKDF)
		}
	}

	return options
}
