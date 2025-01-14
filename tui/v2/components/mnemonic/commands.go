package mnemonic

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/viwet/GoBIP39/words"
	"github.com/viwet/GoDepositCLI/app"
)

type Mnemonic struct {
	mnemonic []string
	list     words.List

	err error
}

func newMnemonicMessage(mnemonic []string, list words.List, err error) *Mnemonic {
	return &Mnemonic{
		mnemonic: mnemonic,
		list:     list,
		err:      err,
	}
}

func generateMnemonic(state *app.State[app.DepositConfig]) tea.Cmd {
	return func() tea.Msg {
		return newMnemonicMessage(app.GenerateMnemonic(state))
	}
}
