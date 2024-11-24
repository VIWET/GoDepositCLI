package config

import (
	"encoding/json"

	"github.com/viwet/GoDepositCLI/helpers"
)

const (
	// 1 Ether = 1.000.000.000 Gwei
	GweiPerEther = 1e9

	// Prefix used in BLS withdrawal credentials
	BLSWithdrawalPrefix = 0x00

	// Prefix used in execution address withdrawal credentials
	ExecutionAddressWithdrawalPrefix = 0x01

	// Execution address length
	ExecutionAddressLength = 20

	// Sha256 hash length
	HashLength = 32

	// Fork version length
	ForkVersionLength = 4

	GweiSuffix = "GWEI"
)

// Chain config
type ChainConfig struct {
	Name                  string
	GenesisForkVersion    []byte
	GenesisValidatorsRoot []byte
}

type ChainConfigJSON struct {
	Name                  string      `json:"name"`
	GenesisForkVersion    helpers.Hex `json:"genesis_fork_version"`
	GenesisValidatorsRoot helpers.Hex `json:"genesis_validators_root"`
}

func (cfg *ChainConfig) UnmarshalJSON(data []byte) error {
	var decoded ChainConfigJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	cfg.Name = decoded.Name
	cfg.GenesisForkVersion = decoded.GenesisForkVersion
	cfg.GenesisValidatorsRoot = decoded.GenesisValidatorsRoot

	return nil
}
