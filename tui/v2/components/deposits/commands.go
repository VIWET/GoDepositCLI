package deposits

import (
	"context"
	"encoding/hex"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/types"
	keystore "github.com/viwet/GoKeystoreV4"
)

type result struct {
	deposits  []*types.Deposit
	keystores []*keystore.Keystore
	err       error
}

func newResult(deposits []*types.Deposit, keystores []*keystore.Keystore, err error) result {
	return result{
		deposits:  deposits,
		keystores: keystores,
		err:       err,
	}
}

type deposit struct {
	index     uint32
	publicKey string
}

func newDeposit(index uint32, publicKey []byte) deposit {
	return deposit{
		index:     index,
		publicKey: "0x" + hex.EncodeToString(publicKey),
	}
}

func RunEngine(ctx context.Context, state *app.State[app.DepositConfig]) (tea.Cmd, <-chan deposit) {
	ticks := make(chan deposit, state.Config().EngineWorkers)
	engine := app.NewDepositEngine(state).OnDeposit(func(d *app.Deposit) error {
		index, deposit, _ := d.Unwrap()
		ticks <- newDeposit(index, deposit.PublicKey)
		return nil
	})

	return func() tea.Msg {
		defer close(ticks)
		return newResult(engine.Generate(ctx))
	}, ticks
}

func WaitDeposit(ticks <-chan deposit) tea.Cmd {
	return func() tea.Msg {
		return <-ticks
	}
}
