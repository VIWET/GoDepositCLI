package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoBIP39/words"
	bls12381 "github.com/viwet/GoBLS12381"
	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/types"
	keystore "github.com/viwet/GoKeystoreV4"
)

// GenerateDeposits generates deposit data and keystores from given menmonic and saves them into files
func GenerateDeposits(ctx *cli.Context, mnemonic []string, list words.List) error {
	config, err := NewDepositConfigFromCLI(ctx)
	if err != nil {
		return err
	}

	seed, err := bip39.ExtractSeed(mnemonic, list, "")
	if err != nil {
		return err
	}

	cryptoOptions, err := config.CryptoOptions()
	if err != nil {
		return err
	}

	var (
		deposits  = make([]*types.Deposit, 0, config.Number)
		keystores = make([]*keystore.Keystore, 0, config.Number)

		chain = config.ChainConfig

		from, to uint32 = config.StartIndex, config.StartIndex + config.Number
	)

	password, err := ReadPassword(ctx)
	if err != nil {
		return err
	}

	for index := from; index < to; index++ {
		var (
			signingKeyPath    = bls12381.SigningKeyPath(index)
			withdrawalKeyPath = bls12381.WithdrawalKeyPath(index)
		)

		signingKey, err := deriveSecretKey(seed, signingKeyPath)
		if err != nil {
			return err
		}
		withdrawalKey, err := deriveSecretKey(seed, withdrawalKeyPath)
		if err != nil {
			return err
		}

		deposit, err := types.NewDeposit(
			signingKey,
			withdrawalKey,
			chain,
			config.DepositOptions(index)...,
		)
		if err != nil {
			return err
		}

		keystore, err := keystore.Encrypt(
			signingKey.Marshal(),
			password,
			signingKeyPath,
			keystore.WithCrypto(cryptoOptions...),
			keystore.WithPublicKey(signingKey.PublicKey().Marshal()),
		)
		if err != nil {
			return err
		}

		deposits = append(deposits, deposit)
		keystores = append(keystores, keystore)
	}

	return SaveData(deposits, keystores, config.Directory)
}

// SaveData saves deposits and keystores in given directory
func SaveData(deposits []*types.Deposit, keystores []*keystore.Keystore, directory string) error {
	if err := ensureDirectoryExist(directory); err != nil {
		return err
	}

	if err := SaveDeposits(deposits, directory); err != nil {
		return err
	}

	if err := SaveKeystores(keystores, directory); err != nil {
		return err
	}

	return nil
}

// SaveDeposits save all deposits into file in given directory
func SaveDeposits(deposits []*types.Deposit, directory string) error {
	filepath := path.Join(directory, fmt.Sprintf("deposit_data-%d.json", time.Now().Unix()))
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(deposits); err != nil {
		return err
	}

	return nil
}

// SaveKeystores save each keystore into file in given directory
func SaveKeystores(keystores []*keystore.Keystore, directory string) error {
	for _, keystore := range keystores {
		if err := saveKeystore(keystore, directory); err != nil {
			return err
		}
	}

	return nil
}

func saveKeystore(keystore *keystore.Keystore, directory string) error {
	pathSuffix := strings.TrimPrefix(strings.ReplaceAll(keystore.Path, "/", "_"), "m")
	filepath := path.Join(directory, fmt.Sprintf("keystore%s.json", pathSuffix))
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(keystore); err != nil {
		return err
	}

	return nil
}

func deriveSecretKey(seed []byte, path string) (bls.SecretKey, error) {
	keyBig, err := bls12381.DeriveKey(seed, path)
	if err != nil {
		return nil, err
	}

	return bls.UnmarshalSecretKey(bip39.PadLeftToBitlen(keyBig.Bytes(), 256))
}
