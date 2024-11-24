package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoBIP39/words"
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

// ReadMnemonic from flags or stdin
func ReadMnemonic(ctx *cli.Context) ([]string, words.List, error) {
	config, err := NewMnemonicConfigFromCLI(ctx)
	if err != nil {
		return nil, nil, err
	}

	list, err := LanguageToWordList(config.Language)
	if err != nil {
		return nil, nil, err
	}

	var mnemonic []string
	if ctx.IsSet(MnemonicFlag.Name) {
		mnemonic = bip39.SplitMnemonic(strings.TrimSpace(ctx.String(MnemonicFlag.Name)))
	} else {
		var err error
		mnemonic, err = ScanMnemonic()
		if err != nil {
			return nil, nil, err
		}
	}

	if err := bip39.ValidateMnemonic(mnemonic, list); err != nil {
		return nil, nil, err
	}

	return mnemonic, list, nil
}

// ScanMnemonic from stdin
func ScanMnemonic() ([]string, error) {
	fmt.Println("Enter your mnemonic")
	mnemonic, err := readMnemonic()
	if err != nil {
		return nil, err
	}

	mnemonic = strings.TrimSpace(mnemonic)

	return bip39.SplitMnemonic(mnemonic), nil
}

func readMnemonic() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}

// ShowMnemonic writes mnemonic into stdout
func ShowMnemonic(mnemonic []string) {
	fmt.Println(strings.Join(mnemonic, " "))
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
