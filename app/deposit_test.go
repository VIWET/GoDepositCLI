package app

import (
	"bytes"
	"encoding/hex"
	"testing"

	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoBIP39/words"
	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/signing"
	keystore "github.com/viwet/GoKeystoreV4"
)

func Test_GenerateDeposits(t *testing.T) {
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
		cfg = &DepositConfig{
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
			WithdrawalAddresses: &IndexedConfigWithDefault[Address]{
				Default: dead0,
				IndexedConfig: IndexedConfig[Address]{
					Config: map[uint32]Address{2: dead2},
				},
			},
			KeystoreKeyDerivationFunction: "scrypt",
		}

		mnemonic = bip39.SplitMnemonic("abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art")
		list     = words.English
		password = "𝔱𝔢𝔰𝔱𝔭𝔞𝔰𝔰𝔴𝔬𝔯𝔡🔑"

		from, to uint32 = cfg.StartIndex, cfg.StartIndex + cfg.Number
	)

	if err := EnsureDepositConfigIsValid(cfg); err != nil {
		t.Fatal(err)
	}

	deposits, keystores, err := GenerateDeposits(cfg, mnemonic, list, password)
	if err != nil {
		t.Fatal(err)
	}

	if len(deposits) != len(keystores) {
		t.Fatalf("generated different number of keystores and deposits - deposits: %d, keystores: %d", len(deposits), len(keystores))
	}

	for i := from; i < to; i++ {
		var (
			d  = deposits[i]
			ks = keystores[i]
		)

		signingKeyBytes, err := keystore.Decrypt(ks, password)
		if err != nil {
			t.Fatal(err)
		}

		signingKey, err := bls.UnmarshalSecretKey(signingKeyBytes)
		if err != nil {
			t.Fatal(err)
		}

		publicKey := signingKey.PublicKey()
		publicKeyBytes := publicKey.Marshal()

		if !bytes.Equal(d.PublicKey, publicKeyBytes) {
			t.Fatalf("deposit public key and extracted public key are not equal\nDeposit:   0x%x\nExtracted: 0x%x", d.PublicKey, publicKeyBytes)
		}

		if !bytes.Equal(d.PublicKey, ks.PublicKey) {
			t.Fatalf("deposit public key and keystore public key are not equal\nDeposit:  0x%x\nKeystore: 0x%x", d.PublicKey, ks.PublicKey)
		}

		signature, err := bls.UnmarshalSignature(d.Signature)
		if err != nil {
			t.Fatal(err)
		}
		domain, err := signing.DepositDomain(cfg.ChainConfig.GenesisForkVersion)
		if err != nil {
			t.Fatal(err)
		}
		data := signing.SigningData{Root: d.DepositMessageRoot, Domain: domain}
		root, err := data.HashTreeRoot()
		if err != nil {
			t.Fatal(err)
		}

		if !signature.Verify(publicKey, root[:]) {
			t.Fatal("invalid signature")
		}

		if amount := cfg.Amounts.Get(i); amount.Gwei() != d.Amount {
			t.Fatalf("invalid amount - want: %d, got: %d", amount, d.Amount)
		}

		withdrawalAddress := cfg.WithdrawalAddresses.Get(i)
		if d.WithdrawalCredentials[0] != 0x01 && !bytes.Equal(d.WithdrawalCredentials[12:], withdrawalAddress[:]) {
			t.Fatal("invalid withdrawal address")
		}

	}
}