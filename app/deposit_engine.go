package app

import (
	"context"
	"runtime"
	"sync"

	bip39 "github.com/viwet/GoBIP39"
	"github.com/viwet/GoDepositCLI/helpers"
	"github.com/viwet/GoDepositCLI/types"
	keystore "github.com/viwet/GoKeystoreV4"
)

// Deposit...
type Deposit struct {
	index    uint32
	deposit  *types.Deposit
	keystore *keystore.Keystore
}

func newDeposit(index uint32, deposit *types.Deposit, keystore *keystore.Keystore) *Deposit {
	return &Deposit{
		index:    index,
		deposit:  deposit,
		keystore: keystore,
	}
}

func (d *Deposit) Unwrap() (uint32, *types.Deposit, *keystore.Keystore) {
	return d.index, d.deposit, d.keystore
}

// DepositEngine...
type DepositEngine struct {
	*State[DepositConfig]

	onDepositFunc func(*Deposit) error
}

// NewDepositEngine...
func NewDepositEngine(state *State[DepositConfig]) *DepositEngine {
	return &DepositEngine{
		State: state,
	}
}

// OnDeposit sets the function which will be called once Deposit and Keystore is generated
func (e *DepositEngine) OnDeposit(onDeposit func(*Deposit) error) *DepositEngine {
	e.onDepositFunc = onDeposit
	return e
}

// Generate generates all deposits and keystores according to the config concurrently
func (e *DepositEngine) Generate(ctx context.Context) ([]*types.Deposit, []*keystore.Keystore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		cfg = e.cfg

		from, to uint32 = cfg.StartIndex, cfg.StartIndex + cfg.Number
	)

	results, err := e.generateDeposits(ctx, indices(ctx, from, to))
	if err != nil {
		return nil, nil, err
	}

	return e.collectDeposits(ctx, results)
}

func (e *DepositEngine) onDeposit(d *Deposit) error {
	if e.onDepositFunc != nil {
		return e.onDepositFunc(d)
	}
	return nil
}

func (e *DepositEngine) generateDeposits(ctx context.Context, indices <-chan uint32) (<-chan helpers.Result[*Deposit], error) {
	var (
		cfg      = e.cfg
		mnemonic = e.mnemonic
		list     = e.list
		password = e.password

		workers = max(1, min(cfg.EngineWorkers, runtime.NumCPU()))
		results = make(chan helpers.Result[*Deposit], workers)
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

				d, k, err := generateDeposit(cfg, seed, i, password)
				if err != nil {
					results <- helpers.Error[*Deposit](err)
					continue
				}

				results <- helpers.Ok[*Deposit](newDeposit(i, d, k))
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

func (e *DepositEngine) collectDeposits(ctx context.Context, results <-chan helpers.Result[*Deposit]) ([]*types.Deposit, []*keystore.Keystore, error) {
	var (
		cfg = e.cfg

		deposits  = make([]*types.Deposit, cfg.Number)
		keystores = make([]*keystore.Keystore, cfg.Number)
	)

	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case result, ok := <-results:
			if !ok {
				return deposits, keystores, nil
			}

			deposit, err := result.Unwrap()
			if err != nil {
				return nil, nil, err
			}

			if err := e.onDeposit(deposit); err != nil {
				return nil, nil, err
			}

			i, d, k := deposit.Unwrap()
			index := int(i - cfg.StartIndex)
			deposits[index] = d
			keystores[index] = k
		}
	}
}

func indices(ctx context.Context, from, to uint32) <-chan uint32 {
	indices := make(chan uint32)
	go func() {
		defer close(indices)
		for index := from; index < to; index++ {
			select {
			case <-ctx.Done():
				return
			case indices <- index:
			}
		}
	}()

	return indices
}
