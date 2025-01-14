package mnemonic

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/tui/v2/components/menu"
)

var bitlens = [...]uint{
	128,
	160,
	192,
	224,
	256,
}

func NewBitlen(ctx *cli.Context, state *app.State[app.DepositConfig]) (tea.Model, tea.Cmd) {
	return menu.New("Bitlen", generateBitlenOptions(ctx, state)...), nil
}

func generateBitlenOptions(ctx *cli.Context, state *app.State[app.DepositConfig]) []menu.Option {
	options := make([]menu.Option, len(bitlens))
	for i, bitlen := range bitlens {
		title := fmt.Sprintf("%d (%d words)", bitlen, (bitlen+(bitlen/32))/11)
		options[i] = menu.NewOption(title, func() (tea.Model, tea.Cmd) {
			state.WithMnemonic(nil, nil)
			state.Config().MnemonicConfig.Bitlen = bitlen
			return New(ctx, state)
		})
	}

	return options
}
