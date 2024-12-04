package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/viwet/GoDepositCLI/config"
	keystore "github.com/viwet/GoKeystoreV4"
)

// Config validation errors
var (
	ErrNumberIsZero      = errors.New("number cannot be 0")
	ErrNoOutputDirectory = errors.New("output directory must be provided")

	ErrInvalidGenesisForkVersion    = fmt.Errorf("invalid genesis fork version length, must be 4 bytes")
	ErrInvalidGenesisValidatorsRoot = fmt.Errorf("invalid genesis validators root length, must be 32 bytes")

	ErrInvalidMnemonicLanguage = fmt.Errorf("invalid language (only %s allowed)", strings.Join(allowedLanguagesNames[:], ", "))
	ErrInvalidMnemonicBitlen   = fmt.Errorf("invalid bitlen (only 128, 160, 192, 244, 256 allowed)")
)

// DepositConfig validation errors
var (
	ErrInvalidKDF = fmt.Errorf(
		"invalid key derivation function (only %s allowed)",
		strings.Join([]string{keystore.ScryptName, keystore.PBKDF2Name}, ", "),
	)

	ErrInvalidAmount = fmt.Errorf(
		"invalid amount (should be between %d and %d and divisible by %d)",
		config.MinDepositAmount,
		config.MaxDepositAmount,
		uint64(config.GweiPerEther),
	)
)
