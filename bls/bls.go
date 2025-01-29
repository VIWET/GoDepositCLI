package bls

var (
	_ SecretKey = (&secretKey{})
	_ PublicKey = (&publicKey{})
	_ Signature = (&signature{})
)

var dst = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

// BLS Secret key
type SecretKey interface {
	PublicKey() PublicKey
	Marshal() []byte
	Sign(message []byte) Signature
}

// BLS Public key
type PublicKey interface {
	Marshal() []byte
}

// BLS Signature
type Signature interface {
	Verify(publicKey PublicKey, message []byte) bool
	Marshal() []byte
}
