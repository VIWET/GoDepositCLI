//go:build !ethereum

package app

// DepositConfig stores all deposit generation related data
type DepositConfig struct {
	*Config

	Amounts             *IndexedConfigWithDefault[uint64]  `json:"amounts,omitempty"`
	WithdrawalAddresses *IndexedConfigWithDefault[Address] `json:"withdrawal_addresses,omitempty"`

	ContractAddresses *IndexedConfig[Address] `json:"contract_addresses,omitempty"`

	KeystoreKeyDerivationFunction string `json:"kdf,omitempty"`
}
