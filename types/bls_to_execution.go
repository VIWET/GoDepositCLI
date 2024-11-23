package types

import (
	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
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
	Message   BLSToExecution
	Signature []byte
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
