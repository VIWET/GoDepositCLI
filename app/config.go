package app

import "github.com/viwet/GoDepositCLI/config"

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
