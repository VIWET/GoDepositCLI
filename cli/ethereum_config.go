//go:build ethereum

package cli

// DepositConfig stores all deposit related params
type DepositConfig struct {
	StartIndex uint32 `json:"start_index"`
	Number     uint32 `json:"number"`

	ChainConfig *config.ChainConfig `json:"chain_config,omitempty"`

	Amounts             *IndexedConfigWithDefault[uint64]  `json:"amount,omitempty"`
	WithdrawalAddresses *IndexedConfigWithDefault[Address] `json:"withdrawal_addresses,omitempty"`

	Directory string `json:"directory"`

	KeystoreKeyDerivationFunction string `json:"kdf,omitempty"`
}
