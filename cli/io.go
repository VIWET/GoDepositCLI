package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoDepositCLI/app"
)

func ShowMnemonic(state *app.State[app.DepositConfig]) {
	fmt.Println(strings.Join(state.Mnemonic(), " "))
}

func ReadPassword(cmd *cli.Command) (string, error) {
	if !cmd.IsSet(PasswordFlag.Name) {
		return "", errors.New("cannot read password in non-interactive mode, --password flag must be set")
	}

	return cmd.String(PasswordFlag.Name), nil
}

func ReadMnemonic(cmd *cli.Command) ([]string, error) {
	if !cmd.IsSet(MnemonicFlag.Name) {
		return nil, errors.New("cannot read mnemonic in non-interactive mode, --mnemonic flag must be set")
	}

	return bip39.SplitMnemonic(strings.TrimSpace(cmd.String(MnemonicFlag.Name))), nil
}
