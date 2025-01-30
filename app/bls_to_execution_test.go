package app

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoBIP39/words"
	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/signing"
)

func Test_GenerateBLSToExecutionMessages(t *testing.T) {
	tempdir := t.TempDir()

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

	var (
		dcfg = &DepositConfig{
			Config: &Config{
				StartIndex: 0,
				Number:     3,
				ChainConfig: &config.ChainConfig{
					Name:                  "devnet",
					GenesisForkVersion:    genesisForkVersion,
					GenesisValidatorsRoot: genesisValidatorsRoot,
				},
				MnemonicConfig: &MnemonicConfig{
					Language: "english",
					Bitlen:   256,
				},
				Directory: tempdir,
			},
			Amounts: &IndexedConfigWithDefault[Amount]{
				Default: Amount(config.MaxDepositAmount / 2),
				IndexedConfig: IndexedConfig[Amount]{
					Config: map[uint32]Amount{
						0: Amount(config.MaxDepositAmount),
						1: Amount(config.MinDepositAmount),
					},
				},
			},
			KeystoreKeyDerivationFunction: "scrypt",
		}

		blscfg = &BLSToExecutionConfig{
			Config: &Config{
				StartIndex: 0,
				Number:     3,
				ChainConfig: &config.ChainConfig{
					Name:                  "devnet",
					GenesisForkVersion:    genesisForkVersion,
					GenesisValidatorsRoot: genesisValidatorsRoot,
				},
				MnemonicConfig: &MnemonicConfig{
					Language: "english",
					Bitlen:   256,
				},
				Directory: tempdir,
			},
			WithdrawalAddresses: &IndexedConfigWithDefault[Address]{
				Default: dead0,
				IndexedConfig: IndexedConfig[Address]{
					Config: map[uint32]Address{
						1: dead1,
						2: dead2,
					},
				},
			},
			ValidatorIndices: &IndexedConfig[uint64]{
				Config: map[uint32]uint64{
					0: 0,
					1: 1,
					2: 2,
				},
			},
		}

		mnemonic = bip39.SplitMnemonic("abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art")
		list     = words.English
		password = "𝔱𝔢𝔰𝔱𝔭𝔞𝔰𝔰𝔴𝔬𝔯𝔡🔑"

		from, to uint32 = dcfg.StartIndex, dcfg.StartIndex + dcfg.Number
	)

	if err := EnsureDepositConfigIsValid(dcfg); err != nil {
		t.Fatal(err)
	}

	if err := EnsureBLSToExecutionConfigIsValid(blscfg); err != nil {
		t.Fatal(err)
	}

	dstate := NewState(dcfg).WithMnemonic(mnemonic, list).WithPassword(password)
	blsstate := NewState(blscfg).WithMnemonic(mnemonic, list)

	deposits, keystores, err := GenerateDeposits(dstate)
	if err != nil {
		t.Fatal(err)
	}

	if len(deposits) != len(keystores) {
		t.Fatalf("generated different number of keystores and deposits - deposits: %d, keystores: %d", len(deposits), len(keystores))
	}

	messages, err := GenerateBLSToExecutionMessages(blsstate)
	if err != nil {
		t.Fatal(err)
	}

	if len(deposits) != len(keystores) {
		t.Fatalf("generated different number of bls to execution and deposits - deposits: %d, bls to execution: %d", len(deposits), len(messages))
	}

	for i := from; i < to; i++ {
		message := messages[i-from]
		deposit := deposits[i-from]
		credentials := sha256.Sum256(message.Message.FromBLSPublicKey)
		if !bytes.Equal(credentials[1:], deposit.WithdrawalCredentials[1:]) {
			t.Fatal("invalid credentials")
		}

		publicKey, err := bls.UnmarshalPublicKey(message.Message.FromBLSPublicKey)
		if err != nil {
			t.Fatal(err)
		}

		signature, err := bls.UnmarshalSignature(message.Signature)
		if err != nil {
			t.Fatal(err)
		}

		domain, err := signing.BLSToExecutionDomain(genesisForkVersion, genesisValidatorsRoot)
		if err != nil {
			t.Fatal(err)
		}

		msgroot, err := message.Message.HashTreeRoot()
		if err != nil {
			t.Fatal(err)
		}

		data := signing.SigningData{Root: msgroot[:], Domain: domain}
		root, err := data.HashTreeRoot()
		if err != nil {
			t.Fatal(err)
		}

		if !signature.Verify(publicKey, root[:]) {
			t.Fatal("invalid signature")
		}

		withdrawalAddress := blscfg.WithdrawalAddresses.Get(i)
		if !bytes.Equal(message.Message.ToExecutionAddress, withdrawalAddress[:]) {
			t.Fatal("invalid withdrawal address")
		}
	}
}
