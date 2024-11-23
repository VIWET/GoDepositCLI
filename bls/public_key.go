package bls

import (
	"errors"

	blst "github.com/supranational/blst/bindings/go"
)

// BLSPublicKeyLength
const BLSPublicKeyLength = 48

var (
	// ErrPublicKeyInvalidLength
	ErrPublicKeyInvalidLength = errors.New("public key length is not equal to 48 bytes")
	// ErrPublicKeyUnmarshal
	ErrPublicKeyUnmarshal = errors.New("cannot unmarshal BLS12-381 public key from bytes")
)

// publicKey
type publicKey struct {
	*blst.P1Affine
}

// UnmarshalPublicKey returns BLS Public key from bytes
func UnmarshalPublicKey(bytes []byte) (PublicKey, error) {
	if len(bytes) != BLSPublicKeyLength {
		return nil, ErrPublicKeyInvalidLength
	}

	publicKey := &publicKey{new(blst.P1Affine).Uncompress(bytes)}
	if publicKey.P1Affine == nil {
		return nil, ErrPublicKeyUnmarshal
	}

	return publicKey, nil
}

// Marshal returns bytes representation of BLS Public key
func (p *publicKey) Marshal() []byte {
	return p.Compress()
}
