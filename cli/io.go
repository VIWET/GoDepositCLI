package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

// ReadPassword from flag or stdin
func ReadPassword(ctx *cli.Context) (string, error) {
	if password := ctx.String(PasswordFlag.Name); password != "" {
		return password, nil
	}

	return ScanPassword()
}

// ScanPassword from stdin
func ScanPassword() (string, error) {
	fmt.Println("Enter new password:")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	fmt.Println("Confirm your password:")
	for {
		confirm, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}

		if bytes.Equal(password, confirm) {
			break
		}
		fmt.Println("Passwords are not equal, try again!")
	}

	return string(password), nil
}

func ensureDirectoryExist(directory string) error {
	info, err := os.Stat(directory)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return os.MkdirAll(directory, 0o700)
		}
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", directory)
	}

	return nil
}
