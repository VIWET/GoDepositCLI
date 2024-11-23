//go:build ethereum

package types

import (
	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
)

//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path deposit.go --objs DepositData,DepositMessage --include ethereum.go --output deposit.ssz.go
//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path ethereum.go --objs DepositMessage --output deposit_message.ssz.go

// DepositMessage contains Ethereum 2.0 validator deposit data
type DepositMessage struct {
	PublicKey             []byte `ssz-size:"48"`
	WithdrawalCredentials []byte `ssz-size:"32"`
	Amount                uint64
}

// DefaultDepositMessage returns DepositMessage with default values
func DefaultDepositMessage(signingKey, withdrawalKey bls.SecretKey) DepositMessage {
	return DepositMessage{
		PublicKey:             signingKey.PublicKey().Marshal(),
		WithdrawalCredentials: BLSWithdrawalCredentials(withdrawalKey),
		Amount:                config.MaxDepositAmount,
	}
}
