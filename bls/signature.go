package bls

import (
	"errors"

	blst "github.com/supranational/blst/bindings/go"
)

// BLSSignatureLength
const BLSSignatureLength = 96

var (
	// ErrSignatureInvalidLength
	ErrSignatureInvalidLength = errors.New("signature length is not equal to 96 bytes")
	// ErrSignatureUnmarshal
	ErrSignatureUnmarshal = errors.New("cannot unmarshal BLS12-381 signature from bytes")
	// ErrSignatureInvalid
	ErrSignatureInvalid = errors.New("signature not in group")
)

// signature
type signature struct {
	*blst.P2Affine
}

// UnmarshalSignature returns BLS Signature from bytes
func UnmarshalSignature(marshalled []byte) (Signature, error) {
	if len(marshalled) != BLSSignatureLength {
		return nil, ErrSignatureInvalidLength
	}

	signature := &signature{new(blst.P2Affine).Uncompress(marshalled)}
	if signature.P2Affine == nil {
		return nil, ErrSignatureUnmarshal
	}

	if !signature.SigValidate(false) {
		return nil, ErrSignatureInvalid
	}

	return signature, nil
}

// Marshal returns bytes representation of BLS Signature
func (s *signature) Marshal() []byte {
	return s.Compress()
}

// Verify verifies data with provided BLS PublicKey
func (s *signature) Verify(key PublicKey, message []byte) bool {
	pk := key.(*publicKey)
	return s.P2Affine.Verify(false, pk.P1Affine, false, message, dst)
}
