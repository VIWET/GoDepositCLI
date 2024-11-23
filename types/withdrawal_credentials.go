package types

import (
	"crypto/sha256"

	"github.com/viwet/GoDepositCLI/bls"
	"github.com/viwet/GoDepositCLI/config"
)

// BLSWithdrawalCredentials returns credentials from BLS SecretKey
func BLSWithdrawalCredentials(withdrawal bls.SecretKey) []byte {
	hash := sha256.Sum256(withdrawal.Marshal())
	credentials := make([]byte, config.HashLength)
	credentials[0] = config.BLSWithdrawalPrefix
	copy(credentials[1:], hash[1:])

	return credentials
}

// ExecutionAddressWithdrawalCredentials returns execution address in withdrawal credentials format
func ExecutionAddressWithdrawalCredentials(address []byte) []byte {
	credentials := make([]byte, config.HashLength)
	credentials[0] = config.ExecutionAddressWithdrawalPrefix
	copy(credentials[12:], address)

	return credentials
}
