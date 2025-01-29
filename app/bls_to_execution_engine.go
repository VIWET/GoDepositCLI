package app

import (
	"context"
	"sync"

	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoDepositCLI/helpers"
	"github.com/viwet/GoDepositCLI/types"
)

// BLSToExecution...
type BLSToExecution struct {
	index   uint32
	message *types.SignedBLSToExecution
}

func newBLSToExecution(index uint32, message *types.SignedBLSToExecution) *BLSToExecution {
	return &BLSToExecution{
		index:   index,
		message: message,
	}
}

func (m *BLSToExecution) Unwrap() (uint32, *types.SignedBLSToExecution) {
	return m.index, m.message
}

// BLSToExecutionEngine...
type BLSToExecutionEngine struct {
	*State[BLSToExecutionConfig]

	onBLSToExecutionFunc func(*BLSToExecution) error
}

// NewBLSToExecutionEngine...
func NewBLSToExecutionEngine(state *State[BLSToExecutionConfig]) *BLSToExecutionEngine {
	return &BLSToExecutionEngine{
		State: state,
	}
}

// OnBLSToExecution sets the function which will be called once BLSToExecution message is generated
func (e *BLSToExecutionEngine) OnBLSToExecution(onBLSToExecution func(*BLSToExecution) error) *BLSToExecutionEngine {
	e.onBLSToExecutionFunc = onBLSToExecution
	return e
}

// GenerateDeposits generates all bls to execution messages according to the config concurrently
func (e *BLSToExecutionEngine) Generate(ctx context.Context) ([]*types.SignedBLSToExecution, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		cfg = e.cfg

		from, to uint32 = cfg.StartIndex, cfg.StartIndex + cfg.Number
	)

	results, err := e.generateMessages(ctx, indices(ctx, from, to))
	if err != nil {
		return nil, err
	}

	return e.collectMessages(ctx, results)
}

//nolint:unused
func (e *BLSToExecutionEngine) onBLSToExecution(m *BLSToExecution) error {
	if e.onBLSToExecutionFunc != nil {
		return e.onBLSToExecutionFunc(m)
	}
	return nil
}

func (e *BLSToExecutionEngine) generateMessages(ctx context.Context, indices <-chan uint32) (<-chan helpers.Result[*BLSToExecution], error) {
	var (
		cfg      = e.cfg
		mnemonic = e.mnemonic
		list     = e.list

		workers = cfg.EngineWorkers
		results = make(chan helpers.Result[*BLSToExecution], workers)
	)

	seed, err := bip39.ExtractSeed(mnemonic, list, "")
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-indices:
				if !ok {
					return
				}

				m, err := generateBLSToExecutionMessage(cfg, seed, i)
				if err != nil {
					results <- helpers.Error[*BLSToExecution](err)
					continue
				}

				results <- helpers.Ok[*BLSToExecution](newBLSToExecution(i, m))
			}
		}
	}

	for range workers {
		wg.Add(1)
		go worker()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results, nil
}

func (e *BLSToExecutionEngine) collectMessages(ctx context.Context, results <-chan helpers.Result[*BLSToExecution]) ([]*types.SignedBLSToExecution, error) {
	var (
		cfg = e.cfg

		messages = make([]*types.SignedBLSToExecution, cfg.Number)
	)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result, ok := <-results:
			if !ok {
				return messages, nil
			}

			message, err := result.Unwrap()
			if err != nil {
				return nil, err
			}

			if err := e.onBLSToExecutionFunc(message); err != nil {
				return nil, err
			}

			i, m := message.Unwrap()
			index := int(i - cfg.StartIndex)
			messages[index] = m
		}
	}
}
