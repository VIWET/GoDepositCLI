//go:build !ethereum

package types

//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path deposit.go --objs DepositData,DepositMessage --include bahamut.go --output deposit.ssz.go
//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path bahamut.go --objs DepositMessage --output deposit_message.ssz.go

// DepositMessage contains Bahamut chain validator deposit data
type DepositMessage struct {
	PublicKey             []byte `ssz-size:"48"`
	WithdrawalCredentials []byte `ssz-size:"32"`
	ContractAddress       []byte `ssz-size:"20"`
	Amount                uint64
}
