package app

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/viwet/GoDepositCLI/config"
	"golang.org/x/crypto/sha3"
)

// Address type alias
type Address [config.ExecutionAddressLength]byte

// FromHex parses hex string
func (a *Address) FromHex(hexstr string) error {
	hexstr = strings.TrimPrefix(hexstr, "0x")
	if _, err := hex.Decode(a[:], []byte(hexstr)); err != nil {
		return err
	}
	return nil
}

// ToHex returns address in hex format
func (a Address) ToHex() string {
	return string(a.toHex())
}

// ToChecksumHex returns address in hex format with EIP-55 address checksum
func (a Address) ToChecksumHex() string {
	buffer := a.toHex()

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(buffer[2:])
	hash := hasher.Sum(nil)

	for i := 2; i < len(buffer); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buffer[i] > '9' && hashByte > 7 {
			buffer[i] -= 32
		}
	}

	return string(buffer[:])
}

func (a *Address) MarshalJSON() ([]byte, error) {
	hexstr := a.ToChecksumHex()
	return json.Marshal(&hexstr)
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var (
		hexstr string
		err    error
	)

	err = json.Unmarshal(data, &hexstr)
	if err != nil {
		return err
	}

	return a.FromHex(hexstr)
}

func (a Address) toHex() []byte {
	return []byte("0x" + hex.EncodeToString(a[:]))
}
