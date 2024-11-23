//go:build ethereum

package types

//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path deposit.go --objs DepositData,DepositMessage --include ethereum.go --output deposit.ssz.go
//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path ethereum.go --objs DepositMessage --output deposit_message.ssz.go

// DepositMessage contains Ethereum 2.0 validator deposit data
type DepositMessage struct {
	PublicKey             []byte `ssz-size:"48"`
	WithdrawalCredentials []byte `ssz-size:"32"`
	Amount                uint64
}
