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
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	if !IsValidAddress(a) {
		return nil, fmt.Errorf(
			"invalid address: expected length %d, but got %d",
			config.ExecutionAddressLength,
			len(a),
		)
	}

	return a, nil
}

// IsValidAmount returns false if address length is not equal to ExecutionAddressLength
func IsValidAddress(address Address) bool {
	return len(address) == config.ExecutionAddressLength
}

// ParseGwei returns amount in Gwei parsed from string
func ParseGwei(amount string) (uint64, error) {
	suffix, modifier := "", uint64(1)
	switch {
	case strings.HasSuffix(amount, config.GweiSuffix):
		suffix = config.GweiSuffix
	case strings.HasSuffix(amount, config.TokenSuffix):
		suffix, modifier = config.TokenSuffix, config.GweiPerEther
	default:
		modifier = config.GweiPerEther
	}

	amount = strings.TrimSuffix(amount, suffix)
	a, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %w", err)
	}

	a *= modifier

	if !IsValidAmount(a) {
		return 0, fmt.Errorf(
			"invalid amount: expected amount between %d and %d and divisible by %d, but got %d",
			config.MinDepositAmount,
			config.MaxDepositAmount,
			config.GweiPerEther,
			a,
		)
	}

	return a, nil
}

// IsValidAmount returns false if amount less than MinDepositAmount or greater than MaxDepositAmount
func IsValidAmount(amount uint64) bool {
	return config.MinDepositAmount <= amount && amount <= config.MaxDepositAmount && amount%config.GweiPerEther != 0
}

// ParseIndex returns key index between `from` and `to`
func ParseIndex(index string, from, to uint32) (uint32, error) {
	i, err := strconv.ParseUint(index, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid index: %w", err)
	}

	if !IsValidIndex(uint32(i), from, to) {
		return 0, fmt.Errorf(
			"invalid index: expected index between %d and %d, but got %d",
			from,
			to,
			i,
		)
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
		return 0, fmt.Errorf("invalid validator index: %w", err)
	}

	return i, nil
}
