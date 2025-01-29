//go:build !ethereum

package app

import (
	"fmt"
	"testing"
)

func Test_ensureContractAddressesConfigIsValid(t *testing.T) {
	tests := []struct {
		name       string
		makeConfig func(t *testing.T) (*IndexedConfig[Address], uint32, uint32)
		checkError func(t *testing.T, err error)
	}{
		{
			name: "valid contract addresses config",
			makeConfig: func(t *testing.T) (*IndexedConfig[Address], uint32, uint32) {
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

				return &IndexedConfig[Address]{
					Config: map[uint32]Address{
						0: dead0,
						1: dead1,
						2: dead2,
					},
				}, 0, 3
			},
			checkError: nil,
		},
		{
			name: "invalid contract addresses config: invalid key index",
			makeConfig: func(t *testing.T) (*IndexedConfig[Address], uint32, uint32) {
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

				return &IndexedConfig[Address]{
					Config: map[uint32]Address{
						0: dead0,
						1: dead1,
						2: dead2,
					},
				}, 0, 2
			},
			checkError: func(t *testing.T, got error) {
				want := fmt.Errorf(
					"invalid contract addresses config: key index should be between %d and %d, but got %d",
					0,
					2,
					2,
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
			name: "invalid contract addresses config: non-unique address",
			makeConfig: func(t *testing.T) (*IndexedConfig[Address], uint32, uint32) {
				var (
					dead0 Address
					dead1 Address
				)

				if err := dead0.FromHex(deadAddress0); err != nil {
					t.Fatal(err)
				}
				if err := dead1.FromHex(deadAddress1); err != nil {
					t.Fatal(err)
				}

				return &IndexedConfig[Address]{
					Config: map[uint32]Address{
						0: dead0,
						1: dead1,
						2: dead1,
					},
				}, 0, 3
			},
			checkError: func(t *testing.T, got error) {
				var dead1 Address
				if err := dead1.FromHex(deadAddress1); err != nil {
					t.Fatal(err)
				}

				want1 := fmt.Errorf(
					"invalid contract addresses config: %d and %d have the same contract address %s",
					2,
					1,
					dead1.ToChecksumHex(),
				)
				want2 := fmt.Errorf(
					"invalid contract addresses config: %d and %d have the same contract address %s",
					1,
					2,
					dead1.ToChecksumHex(),
				)

				if got == nil {
					t.Fatalf("want error %v or %v, but got nil", want1, want2)
				}

				if want1.Error() != got.Error() && want2.Error() != got.Error() {
					t.Fatalf("invalid error\nwant: %v\ngot:  %v", want1, got)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, from, to := test.makeConfig(t)
			err := ensureContractAddressesConfigIsValid(config, from, to)
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
