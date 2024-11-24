package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/urfave/cli/v2"
	bip39 "github.com/viwet/GoBIP39"
	bls12381 "github.com/viwet/GoBLS12381"
	"github.com/viwet/GoDepositCLI/types"
)

// GenerateBLSToExecution is a cli.Action
func GenerateBLSToExecution(ctx *cli.Context) error {
	mnemonic, list, err := ReadMnemonic(ctx)
	if err != nil {
		return err
	}

	config, err := NewBLSToExecutionConfigFromCLI(ctx)
	if err != nil {
		return err
	}

	seed, err := bip39.ExtractSeed(mnemonic, list, "")
	if err != nil {
		return err
	}

	var (
		messages = make([]*types.SignedBLSToExecution, 0, config.Number)

		chain = config.ChainConfig

		from, to uint32 = config.StartIndex, config.StartIndex + config.Number
	)

	for index := from; index < to; index++ {
		withdrawalKeyPath := bls12381.WithdrawalKeyPath(index)

		withdrawalKey, err := deriveSecretKey(seed, withdrawalKeyPath)
		if err != nil {
			return err
		}

		validatorIndex, ok := config.ValidatorIndices.Get(index)
		if !ok {
			return fmt.Errorf("validator index for key %d not found", index)
		}

		message, err := types.NewBLSToExecution(
			withdrawalKey,
			chain,
			validatorIndex,
			config.WithdrawalAddresses.Get(index),
		)
		if err != nil {
			return err
		}

		messages = append(messages, message)
	}

	return SaveBLSToExectuionMessages(messages, config.Directory)
}

// SaveBLSToExectuionMessages save all bls to execution messages into file in given directory
func SaveBLSToExectuionMessages(messages []*types.SignedBLSToExecution, directory string) error {
	if err := ensureDirectoryExist(directory); err != nil {
		return err
	}

	filepath := path.Join(directory, fmt.Sprintf("bls_to_execution-%d.json", time.Now().Unix()))
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(messages); err != nil {
		return err
	}

	return nil
}
