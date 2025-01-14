package password

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const MinPasswordLength = 8

var (
	ErrBalnkPassword = errors.New("password cannot contain only blank spaces")
	ErrShortPassword = errors.New("password must be at least 8 characters")

	ErrNotEqual = errors.New("passwords are not equal")
)

func validatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return ErrBalnkPassword
	}

	if utf8.RuneCountInString(password) < MinPasswordLength {
		return ErrShortPassword
	}

	return nil
}

func validateConfirmation(password, confirmation string) error {
	if password != confirmation {
		return ErrNotEqual
	}

	return nil
}
