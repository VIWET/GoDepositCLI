//go:build !ethereum

package types

import (
	"fmt"

	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
)

//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path deposit.go --objs DepositData,DepositMessage --include bahamut.go --output deposit.ssz.go
//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path bahamut.go --objs DepositMessage --output deposit_message.ssz.go

// DepositMessage contains Bahamut chain validator deposit data
type DepositMessage struct {
	PublicKey             []byte `ssz-size:"48"`
	WithdrawalCredentials []byte `ssz-size:"32"`
	ContractAddress       []byte `ssz-size:"20"`
	Amount                uint64
}

// zeroContract
var zeroContract [config.ExecutionAddressLength]byte

// DefaultDepositMessage returns DepositMessage with default values
func DefaultDepositMessage(signingKey, withdrawalKey bls.SecretKey) DepositMessage {
	return DepositMessage{
		PublicKey:             signingKey.PublicKey().Marshal(),
		WithdrawalCredentials: BLSWithdrawalCredentials(withdrawalKey),
		ContractAddress:       zeroContract[:],
		Amount:                config.MaxDepositAmount,
	}
}

// WithContract sets contract address
func WithContract(address []byte) DepositOption {
	return func(message *DepositMessage) error {
		if len(address) != config.ExecutionAddressLength {
			return fmt.Errorf("invalid contract address length")
		}

		message.ContractAddress = address
		return nil
	}
}
