package bls

import (
	"errors"

	blst "github.com/supranational/blst/bindings/go"
)

// BLSSecretKeyLength
const BLSSecretKeyLength = 32

var (
	// ErrSecretKeyInvalidLength
	ErrSecretKeyInvalidLength = errors.New("secret key length is not equal to 32 bytes")
	// ErrSecretKeyUnmarshal
	ErrSecretKeyUnmarshal = errors.New("cannot unmarshal BLS12-381 secret key from bytes")
)

// secretKey
type secretKey struct {
	*blst.SecretKey
}

// UnmarshalSecretKey returns BLS Secret key from bytes
func UnmarshalSecretKey(bytes []byte) (SecretKey, error) {
	if len(bytes) != BLSSecretKeyLength {
		return nil, ErrSecretKeyInvalidLength
	}

	secretKey := &secretKey{new(blst.SecretKey).Deserialize(bytes)}
	if secretKey.SecretKey == nil {
		return nil, ErrSecretKeyUnmarshal
	}

	return secretKey, nil
}

// Marshal returns bytes representation of BLS Secret key
func (s *secretKey) Marshal() []byte {
	return s.Serialize()
}

// PublicKey returns corresponding BLS Public key
func (s *secretKey) PublicKey() PublicKey {
	return &publicKey{new(blst.P1Affine).From(s.SecretKey)}
}

// Sign signs provided data
func (s *secretKey) Sign(message []byte) Signature {
	return &signature{new(blst.P2Affine).Sign(s.SecretKey, message, dst)}
}
