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

func Test_EnsureDepositConfigIsValid(t *testing.T) {
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
					Amounts: &IndexedConfigWithDefault[Amount]{
						Default: Amount(config.MinDepositAmount),
						IndexedConfig: IndexedConfig[Amount]{
							Config: map[uint32]Amount{
								1: Amount(config.MaxDepositAmount / 2),
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
			name: "invalid config: invalid validators root (nil error because of skipping validation)",
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
			error: nil,
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
			name: "invalid config: mainnet chain config with different validators root (nil error because of skipping validation)",
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
			error: nil,
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
					Amounts: &IndexedConfigWithDefault[Amount]{
						Default: Amount(config.MinDepositAmount - 1),
					},
				}
			},
			error: fmt.Errorf("invalid default amount %d: %w", config.MinDepositAmount-1, ErrInvalidAmount),
		},
		{
			name: "invalid config: invalid key index",
			makeConfig: func(t *testing.T) *DepositConfig {
				return &DepositConfig{
					Amounts: &IndexedConfigWithDefault[Amount]{
						IndexedConfig: IndexedConfig[Amount]{
							Config: map[uint32]Amount{
								1: Amount(config.MaxDepositAmount),
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
					Amounts: &IndexedConfigWithDefault[Amount]{
						IndexedConfig: IndexedConfig[Amount]{
							Config: map[uint32]Amount{
								0: Amount(config.MinDepositAmount - 1),
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
			err := EnsureDepositConfigIsValid(config)
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

func Test_EnsureBLSToExecutionConfigIsValid(t *testing.T) {
	tests := []struct {
		name       string
		makeConfig func(t *testing.T) *BLSToExecutionConfig
		checkError func(t *testing.T, err error)
	}{
		{
			name: "valid config",
			makeConfig: func(t *testing.T) *BLSToExecutionConfig {
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

				return &BLSToExecutionConfig{
					Config: &Config{
						Number: 2,
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
					ValidatorIndices: &IndexedConfig[uint64]{
						Config: map[uint32]uint64{
							0: 0,
							1: 1,
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
				}
			},
			checkError: nil,
		},
		{
			name: "invalid config: invalid withdrawal addresses config",
			makeConfig: func(t *testing.T) *BLSToExecutionConfig {
				var dead0 Address
				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}

				return &BLSToExecutionConfig{
					ValidatorIndices: &IndexedConfig[uint64]{
						Config: map[uint32]uint64{
							0: 0,
						},
					},
					WithdrawalAddresses: &IndexedConfigWithDefault[Address]{
						IndexedConfig: IndexedConfig[Address]{
							Config: map[uint32]Address{
								1: dead0,
							},
						},
					},
				}
			},
			checkError: func(t *testing.T, got error) {
				want := fmt.Errorf(
					"invalid withdrawal addresses config: key index should be between %d and %d, but got %d",
					0,
					1,
					1,
				)

				if got == nil {
					t.Fatalf("want error %v, but got nil", want)
				}

				if want.Error() != got.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", want, got)
				}
			},
		},
		{
			name: "invalid config: invalid validator indices config",
			makeConfig: func(t *testing.T) *BLSToExecutionConfig {
				var dead0 Address
				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}

				return &BLSToExecutionConfig{
					Config: &Config{
						Number: 3,
					},
					ValidatorIndices: &IndexedConfig[uint64]{
						Config: map[uint32]uint64{
							0: 0,
							1: 1,
							2: 2,
						},
					},
					WithdrawalAddresses: &IndexedConfigWithDefault[Address]{
						IndexedConfig: IndexedConfig[Address]{
							Config: map[uint32]Address{
								0: dead0,
								1: dead0,
							},
						},
					},
				}
			},
			checkError: func(t *testing.T, got error) {
				want := fmt.Errorf("no withdrawal addresses for key indices: %v", []uint32{2})

				if got == nil {
					t.Fatalf("want error %v, but got nil", want)
				}

				if want.Error() != got.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", want, got)
				}
			},
		},
		{
			name: "invalid config: invalid withdrawal addresses config",
			makeConfig: func(t *testing.T) *BLSToExecutionConfig {
				var dead0 Address
				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}

				return &BLSToExecutionConfig{
					Config: &Config{
						Number: 3,
					},
					ValidatorIndices: &IndexedConfig[uint64]{
						Config: map[uint32]uint64{
							0: 0,
							2: 2,
						},
					},
					WithdrawalAddresses: &IndexedConfigWithDefault[Address]{
						Default: dead0,
					},
				}
			},
			checkError: func(t *testing.T, got error) {
				want := fmt.Errorf("no validator indices for key indices: %v", []uint32{1})

				if got == nil {
					t.Fatalf("want error %v, but got nil", want)
				}

				if want.Error() != got.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", want, got)
				}
			},
		},
		{
			name: "invalid config: no withdrawal addresses config",
			makeConfig: func(t *testing.T) *BLSToExecutionConfig {
				var dead0 Address
				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}

				return &BLSToExecutionConfig{
					ValidatorIndices: &IndexedConfig[uint64]{
						Config: map[uint32]uint64{0: 0},
					},
				}
			},
			checkError: func(t *testing.T, got error) {
				want := ErrNoWithdrawalAddresses

				if got == nil {
					t.Fatalf("want error %v, but got nil", want)
				}

				if want.Error() != got.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", want, got)
				}
			},
		},
		{
			name: "invalid config: no validator indices config",
			makeConfig: func(t *testing.T) *BLSToExecutionConfig {
				var dead0 Address
				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}

				return &BLSToExecutionConfig{
					WithdrawalAddresses: &IndexedConfigWithDefault[Address]{
						Default: dead0,
					},
				}
			},
			checkError: func(t *testing.T, got error) {
				want := ErrNoValidatorIndices

				if got == nil {
					t.Fatalf("want error %v, but got nil", want)
				}

				if want.Error() != got.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", want, got)
				}
			},
		},
		{
			name: "invalid config: no validator indices config",
			makeConfig: func(t *testing.T) *BLSToExecutionConfig {
				var dead0 Address
				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}

				return &BLSToExecutionConfig{
					Config: &Config{
						Number: 3,
					},
					ValidatorIndices: &IndexedConfig[uint64]{
						Config: map[uint32]uint64{
							0: 0,
							1: 0,
							2: 2,
						},
					},
					WithdrawalAddresses: &IndexedConfigWithDefault[Address]{
						Default: dead0,
					},
				}
			},
			checkError: func(t *testing.T, got error) {
				var dead1 Address
				if err := dead1.FromHex(deadAddress1); err != nil {
					t.Fatal(err)
				}

				want1 := fmt.Errorf(
					"invalid validator indices config: %d and %d have the same validator index %d",
					0,
					1,
					0,
				)
				want2 := fmt.Errorf(
					"invalid validator indices config: %d and %d have the same validator index %d",
					1,
					0,
					0,
				)

				if got == nil {
					t.Fatalf("want error %v or %v, but got nil", want1, want2)
				}

				if want1.Error() != got.Error() && want2.Error() != got.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", want1, got)
				}
			},
		},
		{
			name: "invalid config: invalid validators root",
			makeConfig: func(t *testing.T) *BLSToExecutionConfig {
				var dead0 Address
				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}
				genesisForkVersion, err := hex.DecodeString(deadGenesisForkVersion)
				if err != nil {
					t.Fatal(err)
				}
				invalidGenesisValidatorsRoot, err := hex.DecodeString(deadGenesisValidatorsRoot + deadGenesisValidatorsRoot)
				if err != nil {
					t.Fatal(err)
				}

				return &BLSToExecutionConfig{
					Config: &Config{
						Number: 1,
						ChainConfig: &config.ChainConfig{
							Name:                  "devnet",
							GenesisForkVersion:    genesisForkVersion,
							GenesisValidatorsRoot: invalidGenesisValidatorsRoot,
						},
					},
					ValidatorIndices:    &IndexedConfig[uint64]{Config: map[uint32]uint64{0: 0}},
					WithdrawalAddresses: &IndexedConfigWithDefault[Address]{Default: dead0},
				}
			},
			checkError: func(t *testing.T, got error) {
				want := fmt.Errorf("invalid chain config: %w", ErrInvalidGenesisValidatorsRoot)
				if got == nil {
					t.Fatalf("want error %v, but got nil", want)
				}

				if want.Error() != got.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", want, got)
				}
			},
		},
		{
			name: "invalid config: mainnet chain config with different validators root",
			makeConfig: func(t *testing.T) *BLSToExecutionConfig {
				var dead0 Address
				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}
				genesisValidatorsRoot, err := hex.DecodeString(deadGenesisValidatorsRoot)
				if err != nil {
					t.Fatal(err)
				}

				return &BLSToExecutionConfig{
					Config: &Config{
						Number: 1,
						ChainConfig: &config.ChainConfig{
							Name:                  "mainnet",
							GenesisValidatorsRoot: genesisValidatorsRoot,
						},
					},
					ValidatorIndices:    &IndexedConfig[uint64]{Config: map[uint32]uint64{0: 0}},
					WithdrawalAddresses: &IndexedConfigWithDefault[Address]{Default: dead0},
				}
			},
			checkError: func(t *testing.T, got error) {
				want := fmt.Errorf(
					"invalid chain config: %w",
					fmt.Errorf(
						"different genesis validators root on %s - want: 0x%s, got: 0x%s",
						"mainnet",
						config.MainnetGenesisValidatorsRoot,
						deadGenesisValidatorsRoot,
					),
				)

				if got == nil {
					t.Fatalf("want error %v, but got nil", want)
				}

				if want.Error() != got.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", want, got)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := test.makeConfig(t)
			err := EnsureBLSToExecutionConfigIsValid(config)
			if err != nil {
				if test.checkError == nil {
					t.Fatalf("want no error, but got %v", err)
				}
				test.checkError(t, err)
			} else {
				if test.checkError != nil {
					t.Fatal("want error, but got nil")
				}
			}
		})
	}
}
