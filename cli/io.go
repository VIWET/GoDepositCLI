package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui"
	"github.com/viwet/GoDepositCLI/tui/components/mnemonic"
	"github.com/viwet/GoDepositCLI/tui/components/password"
	"github.com/viwet/GoDepositCLI/types"
	keystore "github.com/viwet/GoKeystoreV4"
)

// TODO(viwet): merge IO operations with TUI in single function

func ShowMnemonic(ctx *cli.Context, state *app.State[app.DepositConfig]) error {
	if ctx.Bool(NonInteractiveFlag.Name) {
		fmt.Println(strings.Join(state.Mnemonic(), " "))
		return nil
	}

	return tui.Run(mnemonic.New(state))
}

func ReadPassword(ctx *cli.Context) (string, error) {
	if ctx.IsSet(PasswordFlag.Name) {
		return ctx.String(PasswordFlag.Name), nil
	} else if ctx.Bool(NonInteractiveFlag.Name) {
		return "", errors.New("cannot read password in non-interactive mode")
	}

	password := password.New()
	if err := tui.Run(password); err != nil {
		return "", err
	}

	return password.Value(), nil
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

func saveBLSToExecution(messages []*types.SignedBLSToExecution, dir string) error {
	filePath := blsToExecutionPath(dir)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, FilePermission)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(messages)
}

func depositsFilePath(dir string) string {
	return path.Join(dir, fmt.Sprintf("deposit_data-%d.json", time.Now().Unix()))
}

func keystoreFilePath(dir string, keyPath string) string {
	pathSuffix := strings.TrimPrefix(strings.ReplaceAll(keyPath, "/", "_"), "m")
	return path.Join(dir, fmt.Sprintf("keystore%s.json", pathSuffix))
}

func blsToExecutionPath(dir string) string {
	return path.Join(dir, fmt.Sprintf("bls_to_execution-%d.json", time.Now().Unix()))
}
