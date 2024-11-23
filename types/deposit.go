package types

import (
	"fmt"

	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/helpers"
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

// Type alias
type DepositOption = helpers.Option[*DepositMessage]

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
