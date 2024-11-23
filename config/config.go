package config

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
)

// Chain config
type ChainConfig struct {
	Name                  string
	GenesisForkVersion    []byte
	GenesisValidatorsRoot []byte
}
