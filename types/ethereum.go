//go:build ethereum

package types

import (
	"encoding/json"

	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/helpers"
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

type DepositJSON struct {
	PublicKey             helpers.Hex `json:"pubkey"`
	WithdrawalCredentials helpers.Hex `json:"withdrawal_credentials"`
	Amount                uint64      `json:"amount"`
	Signature             helpers.Hex `json:"signature"`
	DepositMessageRoot    helpers.Hex `json:"deposit_message_root"`
	DepositDataRoot       helpers.Hex `json:"deposit_data_root"`
	ForkVersion           helpers.Hex `json:"fork_version"`
	NetworkName           string      `json:"network_name"`
	DepositCLIVersion     string      `json:"deposit_cli_version"`
}

func (d *Deposit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DepositJSON{
		d.PublicKey,
		d.WithdrawalCredentials,
		d.Amount,
		d.Signature,
		d.DepositMessageRoot,
		d.DepositDataRoot,
		d.ForkVersion,
		d.NetworkName,
		d.DepositCLIVersion,
	})
}
