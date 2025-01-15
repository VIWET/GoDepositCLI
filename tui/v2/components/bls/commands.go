package deposits

import (
	"context"
	"encoding/hex"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/viwet/GoDepositCLI/app"
	"github.com/viwet/GoDepositCLI/io"
	"github.com/viwet/GoDepositCLI/tui/v2"
	"github.com/viwet/GoDepositCLI/types"
)

type result struct {
	messages []*types.SignedBLSToExecution
	err      error
}

func newResult(messages []*types.SignedBLSToExecution, err error) result {
	return result{
		messages: messages,
		err:      err,
	}
}

type blsToExecution struct {
	index     uint32
	publicKey string
}

func newBLSToExecution(index uint32, publicKey []byte) blsToExecution {
	return blsToExecution{
		index:     index,
		publicKey: "0x" + hex.EncodeToString(publicKey),
	}
}

func RunEngine(ctx context.Context, state *app.State[app.BLSToExecutionConfig]) (tea.Cmd, <-chan blsToExecution) {
	ticks := make(chan blsToExecution, state.Config().EngineWorkers)
	engine := app.NewBLSToExecutionEngine(state).OnBLSToExecution(func(d *app.BLSToExecution) error {
		index, message := d.Unwrap()
		ticks <- newBLSToExecution(index, message.Message.FromBLSPublicKey)
		return nil
	})

	return func() tea.Msg {
		defer close(ticks)
		return newResult(engine.Generate(ctx))
	}, ticks
}

func WaitBLSToExecution(ticks <-chan blsToExecution) tea.Cmd {
	return func() tea.Msg {
		return <-ticks
	}
}

func SaveResult(result result, dir string) tea.Cmd {
	return func() tea.Msg {
		if err := io.EnsureDirectoryExist(dir); err != nil {
			return unwrap(tui.QuitWithError(err))
		}

		if err := io.SaveBLSToExecution(result.messages, dir); err != nil {
			return unwrap(tui.QuitWithError(err))
		}

		return unwrap(tui.Quit())
	}
}

func unwrap(cmd tea.Cmd) tea.Msg {
	return cmd()
}
