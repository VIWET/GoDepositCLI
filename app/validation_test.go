package app

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/viwet/GoDepositCLI/config"
)

const (
	deadGenesisForkVersion    = "deaddead"
	deadGenesisValidatorsRoot = "deaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddead"

	deadAddress0 = "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0000"
	deadAddress1 = "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001"
	deadAddress2 = "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0002"
)

func Test_ensureDepositConfigIsValid(t *testing.T) {
	tests := []struct {
		name       string
		makeConfig func(t *testing.T) *DepositConfig
		error      error
	}{
		{
			name: "valid config: full data",
			makeConfig: func(t *testing.T) *DepositConfig {
				genesisForkVersion, err := hex.DecodeString(deadGenesisForkVersion)
				if err != nil {
					t.Fatal(err)
				}
				genesisValidatorsRoot, err := hex.DecodeString(deadGenesisValidatorsRoot)
				if err != nil {
					t.Fatal(err)
				}

				var (
					dead0 Address
					dead1 Address
					dead2 Address
				)

				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}
				if err := dead1.FromHex(deadAddress1); err != nil {
					t.Fatal(err)
				}
				if err := dead2.FromHex(deadAddress2); err != nil {
					t.Fatal(err)
				}

				return &DepositConfig{
					Config: &Config{
						StartIndex: 0,
						Number:     2,
						ChainConfig: &config.ChainConfig{
							Name:                  "devnet",
							GenesisForkVersion:    genesisForkVersion,
							GenesisValidatorsRoot: genesisValidatorsRoot,
						},
						MnemonicConfig: &MnemonicConfig{
							Language: "english",
							Bitlen:   256,
						},
						Directory: "./validators_data",
					},
					Amounts: &IndexedConfigWithDefault[uint64]{
						Default: config.MinDepositAmount,
						IndexedConfig: IndexedConfig[uint64]{
							Config: map[uint32]uint64{
								1: config.MaxDepositAmount / 2,
							},
						},
					},
					WithdrawalAddresses: &IndexedConfigWithDefault[Address]{
						Default: dead0,
						IndexedConfig: IndexedConfig[Address]{
							Config: map[uint32]Address{
								1: dead1,
							},
						},
					},
					ContractAddresses: &IndexedConfig[Address]{
						Config: map[uint32]Address{
							1: dead2,
						},
					},
					KeystoreKeyDerivationFunction: "scrypt",
				}
			},
			error: nil,
		},
		{
			name: "valid config: no data",
			makeConfig: func(t *testing.T) *DepositConfig {
				return &DepositConfig{}
			},
			error: nil,
		},
		{
			name: "invalid config: invalid fork version",
			makeConfig: func(t *testing.T) *DepositConfig {
				invalidGenesisForkVersion, err := hex.DecodeString(deadGenesisForkVersion + deadGenesisForkVersion)
				if err != nil {
					t.Fatal(err)
				}

				return &DepositConfig{
					Config: &Config{
						ChainConfig: &config.ChainConfig{
							Name:               "devnet",
							GenesisForkVersion: invalidGenesisForkVersion,
						},
					},
				}
			},
			error: fmt.Errorf("invalid chain config: %w", ErrInvalidGenesisForkVersion),
		},
		{
			name: "invalid config: invalid validators root",
			makeConfig: func(t *testing.T) *DepositConfig {
				genesisForkVersion, err := hex.DecodeString(deadGenesisForkVersion)
				if err != nil {
					t.Fatal(err)
				}
				invalidGenesisValidatorsRoot, err := hex.DecodeString(deadGenesisValidatorsRoot + deadGenesisValidatorsRoot)
				if err != nil {
					t.Fatal(err)
				}

				return &DepositConfig{
					Config: &Config{
						ChainConfig: &config.ChainConfig{
							Name:                  "devnet",
							GenesisForkVersion:    genesisForkVersion,
							GenesisValidatorsRoot: invalidGenesisValidatorsRoot,
						},
					},
				}
			},
			error: fmt.Errorf("invalid chain config: %w", ErrInvalidGenesisValidatorsRoot),
		},
		{
			name: "invalid config: mainnet chain config with different fork version",
			makeConfig: func(t *testing.T) *DepositConfig {
				genesisForkVersion, err := hex.DecodeString(deadGenesisForkVersion)
				if err != nil {
					t.Fatal(err)
				}

				return &DepositConfig{
					Config: &Config{
						ChainConfig: &config.ChainConfig{
							Name:               "mainnet",
							GenesisForkVersion: genesisForkVersion,
						},
					},
				}
			},
			error: fmt.Errorf(
				"invalid chain config: %w",
				fmt.Errorf(
					"different genesis fork version on mainnet - want: 0x%s, got: 0x%s",
					config.MainnetGenesisForkVersion,
					deadGenesisForkVersion,
				),
			),
		},
		{
			name: "invalid config: mainnet chain config with different validators root",
			makeConfig: func(t *testing.T) *DepositConfig {
				genesisValidatorsRoot, err := hex.DecodeString(deadGenesisValidatorsRoot)
				if err != nil {
					t.Fatal(err)
				}

				return &DepositConfig{
					Config: &Config{
						ChainConfig: &config.ChainConfig{
							Name:                  "mainnet",
							GenesisValidatorsRoot: genesisValidatorsRoot,
						},
					},
				}
			},
			error: fmt.Errorf(
				"invalid chain config: %w",
				fmt.Errorf(
					"different genesis validators root on %s - want: 0x%s, got: 0x%s",
					"mainnet",
					config.MainnetGenesisValidatorsRoot,
					deadGenesisValidatorsRoot,
				),
			),
		},
		{
			name: "invalid config: invalid mnemonic language",
			makeConfig: func(t *testing.T) *DepositConfig {
				return &DepositConfig{
					Config: &Config{
						MnemonicConfig: &MnemonicConfig{
							Language: "russan",
						},
					},
				}
			},
			error: fmt.Errorf("invalid mnemonic config: %w", ErrInvalidMnemonicLanguage),
		},
		{
			name: "invalid config: invalid mnemonic bitlen",
			makeConfig: func(t *testing.T) *DepositConfig {
				return &DepositConfig{
					Config: &Config{
						MnemonicConfig: &MnemonicConfig{
							Bitlen: 255,
						},
					},
				}
			},
			error: fmt.Errorf("invalid mnemonic config: %w", ErrInvalidMnemonicBitlen),
		},
		{
			name: "invalid config: invalid default config",
			makeConfig: func(t *testing.T) *DepositConfig {
				return &DepositConfig{
					Amounts: &IndexedConfigWithDefault[uint64]{
						Default: config.MinDepositAmount - 1,
					},
				}
			},
			error: fmt.Errorf("invalid default amount %d: %w", config.MinDepositAmount-1, ErrInvalidAmount),
		},
		{
			name: "invalid config: invalid key index",
			makeConfig: func(t *testing.T) *DepositConfig {
				return &DepositConfig{
					Amounts: &IndexedConfigWithDefault[uint64]{
						IndexedConfig: IndexedConfig[uint64]{
							Config: map[uint32]uint64{
								1: config.MaxDepositAmount,
							},
						},
					},
				}
			},
			error: fmt.Errorf("invalid amount config: key index should be between %d and %d, but got %d", 0, 1, 1),
		},
		{
			name: "invalid config: invalid amount",
			makeConfig: func(t *testing.T) *DepositConfig {
				return &DepositConfig{
					Amounts: &IndexedConfigWithDefault[uint64]{
						IndexedConfig: IndexedConfig[uint64]{
							Config: map[uint32]uint64{
								0: config.MinDepositAmount - 1,
							},
						},
					},
				}
			},
			error: fmt.Errorf("invalid amount config: invalid amount at index %d (%d): %w",
				0,
				config.MinDepositAmount-1,
				ErrInvalidAmount,
			),
		},
		{
			name: "invalid config: invalid withdrawal addresses config",
			makeConfig: func(t *testing.T) *DepositConfig {
				var dead0 Address
				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}

				return &DepositConfig{
					WithdrawalAddresses: &IndexedConfigWithDefault[Address]{
						IndexedConfig: IndexedConfig[Address]{
							Config: map[uint32]Address{
								1: dead0,
							},
						},
					},
				}
			},
			error: fmt.Errorf("invalid withdrawal addresses config: key index should be between %d and %d, but got %d", 0, 1, 1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := test.makeConfig(t)
			err := ensureDepositConfigIsValid(config)
			if err != nil {
				if err.Error() != test.error.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", test.error, err)
				}
			} else if test.error != nil {
				t.Fatalf("want error %v, but got nil", test.error)
			}
		})
	}
}
