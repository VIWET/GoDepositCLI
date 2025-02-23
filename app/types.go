package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/viwet/GoDepositCLI/config"
	"golang.org/x/crypto/sha3"
)

// Address type
type Address [config.ExecutionAddressLength]byte

var zeroAddress Address

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

// Amount wrapper
type Amount uint64

// ToString returns string representation of amount with suffix
func (a Amount) ToString(suffix string) string {
	var value uint64
	switch suffix {
	case config.GweiSuffix:
		value = a.Gwei()
	case config.TokenSuffix:
		value = a.Ether()
	default:
		value = uint64(a)
	}

	return fmt.Sprintf("%d%s", value, suffix)
}

// Gwei returns amount in Gwei
func (a Amount) Gwei() uint64 {
	return uint64(a)
}

// Ether returns amount in Ether
func (a Amount) Ether() uint64 {
	return uint64(a / config.GweiPerEther)
}

func (a *Amount) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case '"':
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		return a.FromString(str)
	default:
		var amount uint64
		if err := json.Unmarshal(data, &amount); err != nil {
			return err
		}
		*a = Amount(amount)
		return nil
	}
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

// FromString parses value from string (optionally with prefix)
func (a *Amount) FromString(amount string) error {
	suffix, modifier := "", uint64(config.GweiPerEther)
	switch {
	case strings.HasSuffix(amount, config.GweiSuffix):
		suffix, modifier = config.GweiSuffix, 1
	case strings.HasSuffix(amount, config.TokenSuffix):
		suffix, modifier = config.TokenSuffix, config.GweiPerEther
	}

	amount = strings.TrimSuffix(amount, suffix)
	value, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	value *= modifier

	*a = Amount(value)
	return nil
}
