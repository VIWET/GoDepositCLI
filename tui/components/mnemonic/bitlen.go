package mnemonic

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui/components/menu"
)

var bitlens = [...]uint{
	128,
	160,
	192,
	224,
	256,
}

func NewBitlen(ctx context.Context, cmd *cli.Command, state *app.State[app.DepositConfig]) (tea.Model, tea.Cmd) {
	return menu.New("Bitlen", generateBitlenOptions(ctx, cmd, state)...), nil
}

func generateBitlenOptions(ctx context.Context, cmd *cli.Command, state *app.State[app.DepositConfig]) []menu.Option {
	options := make([]menu.Option, len(bitlens))
	for i, bitlen := range bitlens {
		title := fmt.Sprintf("%d (%d words)", bitlen, (bitlen+(bitlen/32))/11)
		options[i] = menu.NewOption(title, func() (tea.Model, tea.Cmd) {
			state.WithMnemonic(nil, nil)
			state.Config().MnemonicConfig.Bitlen = bitlen
			return New(ctx, cmd, state)
		})
	}

	return options
}
