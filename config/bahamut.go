//go:build !ethereum

package config

import "encoding/hex"

const (
	// Minimum allowed deposit in Bahamut chain
	MinDepositAmount uint64 = (1 << 8) * GweiPerEther
	// Maximum allowed deposit in Bahamut chain
	MaxDepositAmount uint64 = (1 << 13) * GweiPerEther
)

// Bahamut Mainnet (Sahara) config
const (
	MainnetName                  = "mainnet"
	MainnetGenesisForkVersion    = "00003341"
	MainnetGenesisValidatorsRoot = "de82e35df831fe1dd9cf93ebc6e09fb042eeb852b448e09d2e3c5ed6b918e038"
)

// MainnetConfig
func MainnetConfig() *ChainConfig {
	genesisForkVersion, err := hex.DecodeString(MainnetGenesisForkVersion)
	if err != nil {
		panic("invalid bahamut MainnetGenesisForkVersion constant")
	}

	genesisValidatorsRoot, err := hex.DecodeString(MainnetGenesisValidatorsRoot)
	if err != nil {
		panic("invalid bahamut MainnetGenesisValidatorsRoot constant")
	}

	return &ChainConfig{
		Name:                  MainnetName,
		GenesisForkVersion:    genesisForkVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
}

// SaharaConfig
func SaharaConfig() *ChainConfig {
	return MainnetConfig()
}

// Horizon testnet config
const (
	HorizonName                  = "horizon"
	HorizonGenesisForkVersion    = "00001934"
	HorizonGenesisValidatorsRoot = "8337cc21ab2223ff7645964cb705d9ea09a3da9c2936ac89209d998e1494fd42"
)

// HorizonConfig
func HorizonConfig() *ChainConfig {
	genesisForkVersion, err := hex.DecodeString(HorizonGenesisForkVersion)
	if err != nil {
		panic("invalid HorizonGenesisForkVersion constant")
	}

	genesisValidatorsRoot, err := hex.DecodeString(HorizonGenesisValidatorsRoot)
	if err != nil {
		panic("invalid HorizonGenesisValidatorsRoot constant")
	}

	return &ChainConfig{
		Name:                  HorizonName,
		GenesisForkVersion:    genesisForkVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
}
