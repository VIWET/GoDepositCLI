//go:build ethereum

package config

import (
	"encoding/hex"
	"strings"
)

const (
	// Minimum allowed deposit in Ethereum
	MinDepositAmount uint64 = (1 << 0) * GweiPerEther
	// Maximum allowed deposit in Ethereum
	MaxDepositAmount uint64 = (1 << 5) * GweiPerEther
	// Token suffix
	TokenSuffix = "ETH"
)

// Ethereum mainnet config
const (
	MainnetName                  = "mainnet"
	MainnetGenesisForkVersion    = "00000000"
	MainnetGenesisValidatorsRoot = "4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95"
)

// MainnetConfig
func MainnetConfig() *ChainConfig {
	genesisForkVersion, err := hex.DecodeString(MainnetGenesisForkVersion)
	if err != nil {
		panic("invalid ethereum MainnetGenesisForkVersion constant")
	}

	genesisValidatorsRoot, err := hex.DecodeString(MainnetGenesisValidatorsRoot)
	if err != nil {
		panic("invalid ethereum MainnetGenesisValidatorsRoot constant")
	}

	return &ChainConfig{
		Name:                  MainnetName,
		GenesisForkVersion:    genesisForkVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
}

// Holesky testnet config
const (
	HoleskyName                  = "holesky"
	HoleskyGenesisForkVersion    = "01017000"
	HoleskyGenesisValidatorsRoot = "9143aa7c615a7f7115e2b6aac319c03529df8242ae705fba9df39b79c59fa8b1"
)

// HoleskyConfig
func HoleskyConfig() *ChainConfig {
	genesisForkVersion, err := hex.DecodeString(HoleskyGenesisForkVersion)
	if err != nil {
		panic("invalid HoleskyGenesisForkVersion constant")
	}

	genesisValidatorsRoot, err := hex.DecodeString(HoleskyGenesisValidatorsRoot)
	if err != nil {
		panic("invalid HoleskyGenesisValidatorsRoot constant")
	}

	return &ChainConfig{
		Name:                  HoleskyName,
		GenesisForkVersion:    genesisForkVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
}

// ConfigByNetworkName
func ConfigByNetworkName(network string) (*ChainConfig, bool) {
	network = strings.ToLower(network)
	switch network {
	case MainnetName:
		return MainnetConfig(), true
	case HoleskyName:
		return HoleskyConfig(), true
	default:
		return nil, false
	}
}
