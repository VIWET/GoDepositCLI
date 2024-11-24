package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/viwet/GoDepositCLI/config"
	"github.com/viwet/GoDepositCLI/helpers"
)

// Address type alias
type Address = helpers.Hex

// ParseAddress from string
func ParseAddress(address string) (Address, error) {
	a, err := hex.DecodeString(strings.TrimPrefix(address, "0x"))
	if err != nil {
		return nil, err
	}

	return a, nil
}

// ParseGwei returns amount in Gwei parsed from string
func ParseGwei(amount string) (uint64, error) {
	switch {
	case strings.HasSuffix(amount, config.GweiSuffix):
		amount, err := parseAmount(amount, config.GweiSuffix)
		if err != nil {
			return 0, err
		}
		return amount, nil
	case strings.HasSuffix(amount, config.TokenSuffix):
		amount, err := parseAmount(amount, config.TokenSuffix)
		if err != nil {
			return 0, err
		}
		return amount * config.GweiPerEther, nil
	default:
		amount, err := parseAmount(amount, "")
		if err != nil {
			return 0, err
		}
		return amount * config.GweiPerEther, nil
	}
}

// parseAmount
func parseAmount(amount string, suffix string) (uint64, error) {
	amount = strings.TrimSuffix(amount, suffix)
	a, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return 0, err
	}

	return a, nil
}

// ParseIndex returns key index between `from` and `to`
func ParseIndex(index string, from, to uint32) (uint32, error) {
	i, err := strconv.ParseUint(index, 10, 32)
	if err != nil {
		return 0, err
	}
	if !IsValidIndex(uint32(i), from, to) {
		return 0, fmt.Errorf("invalid index: expected index between %d and %d, but got %d", from, to, i)
	}

	return uint32(i), nil
}

// IsValidIndex returns false if index less than from or greater than to
func IsValidIndex(index, from, to uint32) bool {
	return from <= index && index < to
}

// ParseValidatorIndex
func ParseValidatorIndex(index string) (uint64, error) {
	i, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}
