package simulators

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

type SimulatedUpkeep struct {
	ID         *big.Int
	EligibleAt []*big.Int
	Performs   map[string]types.PerformLog // performs at block number
}

func (ct *SimulatedContract) GetActiveUpkeepIDs(ctx context.Context) ([]types.UpkeepIdentifier, error) {

	ct.mu.RLock()
	ct.logger.Printf("getting keys at block %s", ct.lastBlock)

	keys := []types.UpkeepIdentifier{}

	// TODO: filter out cancelled upkeeps
	for key := range ct.upkeeps {
		keys = append(keys, types.UpkeepIdentifier(key))
	}
	ct.mu.RUnlock()

	// call to GetState
	err := <-ct.rpc.Call(ctx, "getState")
	if err != nil {
		return nil, err
	}
	// call to GetActiveIDs
	// TODO: batch size is hard coded at 10_000, if the number of keys is more
	// than this, simulate another rpc call
	err = <-ct.rpc.Call(ctx, "getActiveIDs")
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (ct *SimulatedContract) CheckUpkeep(ctx context.Context, keys ...types.UpkeepKey) (types.UpkeepResults, error) {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	var (
		mErr    error
		wg      sync.WaitGroup
		results = make(types.UpkeepResults, len(keys))
	)

	for i, key := range keys {
		wg.Add(1)
		go func(i int, key types.UpkeepKey) {
			defer wg.Done()

			blockKey, upkeepID, err := key.BlockKeyAndUpkeepID()
			if err != nil {
				panic(err.Error())
			}

			block, ok := new(big.Int).SetString(blockKey.String(), 10)
			if !ok {
				mErr = multierr.Append(mErr, fmt.Errorf("block in key not parsable as big int"))
				return
			}

			up, ok := ct.upkeeps[string(upkeepID)]
			if !ok {
				mErr = multierr.Append(mErr, fmt.Errorf("upkeep not registered"))
				return
			}

			var bl [32]byte
			results[i] = types.UpkeepResult{
				Key:     key,
				State:   types.NotEligible,
				GasUsed: big.NewInt(0),
				/*
					FailureReason    uint8
				*/
				PerformData:      []byte{}, // TODO: add perform data
				FastGasWei:       big.NewInt(0),
				LinkNative:       big.NewInt(0),
				CheckBlockNumber: uint32(block.Int64() - 1), // minus 1 because the real contract does this
				CheckBlockHash:   bl,
			}

			// call telemetry after RPC delays have been applied. if a check is cancelled
			// it doesn't count toward telemetry.
			ct.telemetry.CheckKey(key)

			// start at the highest blocks eligible. the first eligible will be a block
			// lower than the current
			for j := len(up.EligibleAt) - 1; j >= 0; j-- {
				e := up.EligibleAt[j]

				if block.Cmp(e) >= 0 {
					results[i].State = types.Eligible

					// check that upkeep has not been recently performed between two
					// points of eligibility
					// is there a log between eligible and block
					var t int64
					diff := new(big.Int).Sub(block, e).Int64()
					for t = 0; t <= diff; t++ {
						c := new(big.Int).Add(e, big.NewInt(t))
						_, ok := up.Performs[c.String()]
						if ok {
							results[i].State = types.NotEligible
							return
						}
					}

					return
				}
			}
		}(i, key)
	}

	wg.Wait()

	if mErr != nil {
		return nil, mErr
	}

	// call to CheckUpkeep
	err := <-ct.rpc.Call(ctx, "checkUpkeep")
	if err != nil {
		return nil, err
	}

	// call to SimulatePerform
	err = <-ct.rpc.Call(ctx, "simulatePerform")
	if err != nil {
		return nil, err
	}

	return results, nil
}
