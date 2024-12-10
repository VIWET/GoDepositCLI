package app

import (
	"fmt"

	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoBIP39/words"
	bls12381 "github.com/viwet/GoBLS12381"
	"github.com/viwet/GoDepositCLI/types"
)

// GenerateDeposits generates all bls to execution messages according to the config
func GenerateBLSToExecutionMessages(cfg *BLSToExecutionConfig, mnemonic []string, list words.List) ([]*types.SignedBLSToExecution, error) {
	seed, err := bip39.ExtractSeed(mnemonic, list, "")
	if err != nil {
		return nil, err
	}

	var (
		messages = make([]*types.SignedBLSToExecution, 0, cfg.Number)

		from, to uint32 = cfg.StartIndex, cfg.StartIndex + cfg.Number
	)

	for index := from; index < to; index++ {
		message, err := generateBLSToExecutionMessage(cfg, seed, index)
		if err != nil {
			return nil, fmt.Errorf("cannot generate bls to execution message for key %d: %w", index, err)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func generateBLSToExecutionMessage(cfg *BLSToExecutionConfig, seed []byte, index uint32) (*types.SignedBLSToExecution, error) {
	withdrawalKeyPath := bls12381.WithdrawalKeyPath(index)
	withdrawalKey, err := deriveSecretKey(seed, withdrawalKeyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot derive withdrawal key: %w", err)
	}

	validatorIndex, ok := cfg.ValidatorIndices.Get(index)
	if !ok {
		return nil, fmt.Errorf("no validator index for key %d", index)
	}

	withdrawalAddress := cfg.WithdrawalAddresses.Get(index)
	message, err := types.NewBLSToExecution(
		withdrawalKey,
		cfg.ChainConfig,
		validatorIndex,
		withdrawalAddress[:],
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create bls to execution message: %w", err)
	}

	return message, nil
}
