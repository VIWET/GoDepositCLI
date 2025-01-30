package types

import (
	"encoding/json"
	"strconv"

	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/helpers"
	"github.com/viwet/GoDepositCLI/signing"
)

//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path bls_to_execution.go --objs BLSToExecution --output bls_to_execution.ssz.go

// BLSToExecution contains data needed for switching from BLS withdrawal credentials to contract address
type BLSToExecution struct {
	ValidatorIndex     uint64
	FromBLSPublicKey   []byte `ssz-size:"48"`
	ToExecutionAddress []byte `ssz-size:"20"`
}

// BLSToExecution contains signed data needed for switching from BLS withdrawal credentials to contract address
type SignedBLSToExecution struct {
	Message   BLSToExecution `json:"message"`
	Signature helpers.Hex    `json:"signature"`
}

// NewBLSToExecution return new signed BLSToExecution message
func NewBLSToExecution(
	withdrawalKey bls.SecretKey,
	config *config.ChainConfig,
	validatorIndex uint64,
	address []byte,
) (*SignedBLSToExecution, error) {
	message := BLSToExecution{
		ValidatorIndex:     validatorIndex,
		FromBLSPublicKey:   withdrawalKey.PublicKey().Marshal(),
		ToExecutionAddress: address,
	}

	messageRoot, err := message.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	domain, err := signing.BLSToExecutionDomain(config.GenesisForkVersion, config.GenesisValidatorsRoot)
	if err != nil {
		return nil, err
	}

	signature, err := signing.SignData(withdrawalKey, messageRoot[:], domain)
	if err != nil {
		return nil, err
	}

	return &SignedBLSToExecution{
		Message:   message,
		Signature: signature.Marshal(),
	}, nil
}

type BLSToExecutionJSON struct {
	ValidatorIndex     string      `json:"validator_index"`
	FromBLSPublicKey   helpers.Hex `json:"from_bls_pubkey"`
	ToExecutionAddress helpers.Hex `json:"to_execution_address"`
}

func (m *BLSToExecution) MarshalJSON() ([]byte, error) {
	return json.Marshal(&BLSToExecutionJSON{
		strconv.FormatUint(m.ValidatorIndex, 10),
		m.FromBLSPublicKey,
		m.ToExecutionAddress,
	})
}
