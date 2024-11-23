package signing

import "github.com/viwet/GoDepositCLI/config"

// BLS Signature domain for beacon deposit verification
var depositDomain = []byte{0x03, 0x00, 0x00, 0x00}

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
