package types

import (
	"fmt"

	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/helpers"
	"github.com/viwet/GoDepositCLI/signing"
	"github.com/viwet/GoDepositCLI/version"
)

// DepositData contains signed validator deposit data
type DepositData struct {
	DepositMessage
	Signature []byte `ssz-size:"96"`
}

// Deposit contains data needed for deposit smart-contract
type Deposit struct {
	DepositData
	DepositMessageRoot []byte
	DepositDataRoot    []byte
	ForkVersion        []byte
	NetworkName        string
	DepositCLIVersion  string
}

type (
	// Type alias
	DepositOption  = helpers.Option[*DepositMessage]
	DepositOptions = helpers.Options[*DepositMessage]
)

// NewDeposit returns new signed Deposit
func NewDeposit(signingKey, withdrawalKey bls.SecretKey, config *config.ChainConfig, opts ...DepositOption) (*Deposit, error) {
	message := DefaultDepositMessage(signingKey, withdrawalKey)
	if err := DepositOptions(opts).Apply(&message); err != nil {
		return nil, err
	}

	messageRoot, err := message.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	domain, err := signing.DepositDomain(config.GenesisForkVersion)
	if err != nil {
		return nil, err
	}

	signature, err := signing.SignData(signingKey, messageRoot[:], domain)
	if err != nil {
		return nil, err
	}

	data := DepositData{message, signature.Marshal()}
	dataRoot, err := data.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	return &Deposit{
		DepositData:        data,
		DepositMessageRoot: messageRoot[:],
		DepositDataRoot:    dataRoot[:],
		ForkVersion:        config.GenesisForkVersion,
		NetworkName:        config.Name,
		DepositCLIVersion:  version.Version(),
	}, nil
}

// WithAmount sets amount
func WithAmount(amount uint64) DepositOption {
	return func(message *DepositMessage) error {
		if amount < config.MinDepositAmount {
			return fmt.Errorf("invalid amount: less than %d", config.MinDepositAmount)
		}
		if amount > config.MaxDepositAmount {
			return fmt.Errorf("invalid amount: greater than %d", config.MaxDepositAmount)
		}
		if amount%config.GweiPerEther != 0 {
			return fmt.Errorf("invalid amount: should be divisible by Ether")
		}

		message.Amount = amount
		return nil
	}
}

// WithWithdrawalAddress sets withdrawal address
func WithWithdrawalAddress(address []byte) DepositOption {
	return func(message *DepositMessage) error {
		if len(address) != config.ExecutionAddressLength {
			return fmt.Errorf("invalid contract address length")
		}

		message.WithdrawalCredentials = ExecutionAddressWithdrawalCredentials(address)
		return nil
	}
}
