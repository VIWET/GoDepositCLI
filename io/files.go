package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/viwet/GoDepositCLI/types"
	keystore "github.com/viwet/GoKeystoreV4"
)

const (
	// -r--------
	FilePermission = 0o400

	// dr-x------
	DirPermission = 0o700
)

func EnsureDirectoryExist(dir string) error {
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

func SaveDeposits(deposits []*types.Deposit, dir string) error {
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

func SaveKeystores(keystores []*keystore.Keystore, dir string) error {
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

func SaveBLSToExecution(messages []*types.SignedBLSToExecution, dir string) error {
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
