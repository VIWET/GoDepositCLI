package signing

import "github.com/viwet/GoDepositCLI/config"

var (
	// BLS Signature domain for beacon deposit verification
	depositDomain = []byte{0x03, 0x00, 0x00, 0x00}
	// BLS Signature domain for beacon bls to execution message verification
	domainBLSToExecution = []byte{0x0A, 0x00, 0x00, 0x00}
)

// DepositDomain returns BLS Signature domain for deposit
func DepositDomain(fork []byte) ([]byte, error) {
	forkData := ForkData{fork, zeroRoot[:]}
	root, err := forkData.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	domain := make([]byte, config.HashLength)
	copy(domain[:config.ForkVersionLength], depositDomain)
	copy(domain[config.ForkVersionLength:], root[:config.HashLength-config.ForkVersionLength])

	return domain, nil
}

// BLSToExecutionDomain returns BLS Signature domain for BLSToExecution message
func BLSToExecutionDomain(fork []byte, validatorsRoot []byte) ([]byte, error) {
	forkData := ForkData{fork, validatorsRoot}
	root, err := forkData.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	domain := make([]byte, 32)
	copy(domain[:config.ForkVersionLength], depositDomain)
	copy(domain[config.ForkVersionLength:], root[:config.HashLength-config.ForkVersionLength])

	return domain, nil
}
