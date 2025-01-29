package signing

import (
	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
)

//go:generate go run -mod=mod github.com/ferranbt/fastssz/sszgen --path signing.go --objs SigningData,ForkData --output signing.ssz.go

// 32 bytes zeroed root
var zeroRoot [config.HashLength]byte

// SigningData container
type SigningData struct {
	Root   []byte `ssz-size:"32"`
	Domain []byte `ssz-size:"32"`
}

// ForkData container
type ForkData struct {
	Version        []byte `ssz-size:"4"`
	ValidatorsRoot []byte `ssz-size:"32"`
}

// SignData messageRoot and domain with provided BLS Secret key
func SignData(secret bls.SecretKey, messageRoot, domain []byte) (bls.Signature, error) {
	data := SigningData{messageRoot, domain}
	root, err := data.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	return secret.Sign(root[:]), nil
}
