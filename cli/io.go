package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui"
	"github.com/viwet/GoDepositCLI/tui/components/mnemonic"
	mnemonicInput "github.com/viwet/GoDepositCLI/tui/components/mnemonic_input"
	"github.com/viwet/GoDepositCLI/tui/components/password"
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
		return "", errors.New("cannot read password in non-interactive mode, --password flag must be set")
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
	} else if ctx.Bool(NonInteractiveFlag.Name) {
		return nil, errors.New("cannot read mnemonic in non-interactive mode, --mnemonic flag must be set")
	}

	mnemonic := mnemonicInput.New()
	if err := tui.Run(mnemonic); err != nil {
		return nil, err
	}

	return bip39.SplitMnemonic(mnemonic.Value()), nil
}
