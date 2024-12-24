package app

import (
	"fmt"

	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoBIP39/words"
	bls12381 "github.com/viwet/GoBLS12381"
	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/types"
	keystore "github.com/viwet/GoKeystoreV4"
)

// GenerateDeposits generates all deposits and keystores according to the config
func GenerateDeposits(state *State[DepositConfig]) ([]*types.Deposit, []*keystore.Keystore, error) {
	var (
		cfg      *DepositConfig = state.cfg
		mnemonic []string       = state.mnemonic
		list     words.List     = state.list
		password string         = state.password
	)

	seed, err := bip39.ExtractSeed(mnemonic, list, "")
	if err != nil {
		return nil, nil, err
	}

	var (
		deposits  = make([]*types.Deposit, 0, cfg.Number)
		keystores = make([]*keystore.Keystore, 0, cfg.Number)

		from, to uint32 = cfg.StartIndex, cfg.StartIndex + cfg.Number
	)

	for index := from; index < to; index++ {
		deposit, keystore, err := generateDeposit(cfg, seed, index, password)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot generate deposit data for key %d: %w", index, err)
		}

		deposits = append(deposits, deposit)
		keystores = append(keystores, keystore)
	}

	return deposits, keystores, nil
}

func generateDeposit(cfg *DepositConfig, seed []byte, index uint32, password string) (*types.Deposit, *keystore.Keystore, error) {
	var (
		cryptoOptions     = newCryptoOptionsFromConfig(cfg)
		signingKeyPath    = bls12381.SigningKeyPath(index)
		withdrawalKeyPath = bls12381.WithdrawalKeyPath(index)
	)

	signingKey, err := deriveSecretKey(seed, signingKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot derive signing key: %w", err)
	}

	withdrawalKey, err := deriveSecretKey(seed, withdrawalKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot derive withdrawal key: %w", err)
	}

	deposit, err := types.NewDeposit(
		signingKey,
		withdrawalKey,
		cfg.ChainConfig,
		newDepositOptionsFromConfig(cfg, index)...,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create deposit: %w", err)
	}

	keystore, err := keystore.Encrypt(
		signingKey.Marshal(),
		password,
		signingKeyPath,
		keystore.WithCrypto(cryptoOptions...),
		keystore.WithPublicKey(signingKey.PublicKey().Marshal()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create keystore: %w", err)
	}

	return deposit, keystore, nil
}

func deriveSecretKey(seed []byte, path string) (bls.SecretKey, error) {
	keyBig, err := bls12381.DeriveKey(seed, path)
	if err != nil {
		return nil, err
	}

	return bls.UnmarshalSecretKey(bip39.PadLeftToBitlen(keyBig.Bytes(), 256))
}
