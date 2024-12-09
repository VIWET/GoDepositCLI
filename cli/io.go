package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoDepositCLI/types"
	keystore "github.com/viwet/GoKeystoreV4"
	"golang.org/x/term"
)

func ShowMnemonic(mnemonic []string) {
	fmt.Println(strings.Join(mnemonic, " "))
}

func ReadPassword(ctx *cli.Context) (string, error) {
	if ctx.IsSet(PasswordFlag.Name) {
		return ctx.String(PasswordFlag.Name), nil
	}

	return scanPassword()
}

func scanPassword() (string, error) {
	for {
		fmt.Println("Enter new password:")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}

		fmt.Println("Confirm your password:")
		confirm, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}

		if bytes.Equal(password, confirm) {
			return string(password), nil
		}
		fmt.Println("Passwords are not equal, try again!")
	}
}

func ReadMnemonic(ctx *cli.Context) ([]string, error) {
	if ctx.IsSet(MnemonicFlag.Name) {
		return bip39.SplitMnemonic(strings.TrimSpace(ctx.String(MnemonicFlag.Name))), nil
	}

	return scanMnemonic()
}

func scanMnemonic() ([]string, error) {
	fmt.Println("Enter your mnemonic")
	reader := bufio.NewReader(os.Stdin)
	mnemonic, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	mnemonic = strings.TrimSpace(mnemonic)
	return bip39.SplitMnemonic(mnemonic), nil
}

const (
	// -r--------
	FilePermission = 0o400

	// dr-x------
	DirPermission = 0o700
)

func ensureDirectoryExist(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return os.MkdirAll(dir, DirPermission)
		}
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	return nil
}

func saveDeposits(deposits []*types.Deposit, dir string) error {
	filePath := depositsFilePath(dir)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, FilePermission)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(deposits)
}

func saveKeystores(keystores []*keystore.Keystore, dir string) error {
	for _, keystore := range keystores {
		if err := saveKeystore(keystore, dir); err != nil {
			return err
		}
	}

	return nil
}

func saveKeystore(keystore *keystore.Keystore, dir string) error {
	filePath := keystoreFilePath(dir, keystore.Path)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, FilePermission)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(keystore)
}

func depositsFilePath(dir string) string {
	return path.Join(dir, fmt.Sprintf("deposit_data-%d.json", time.Now().Unix()))
}

func keystoreFilePath(dir string, keyPath string) string {
	pathSuffix := strings.TrimPrefix(strings.ReplaceAll(keyPath, "/", "_"), "m")
	return path.Join(dir, fmt.Sprintf("keystore%s.json", pathSuffix))
}
