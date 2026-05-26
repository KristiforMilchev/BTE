package implementations

import (
	"bos/interfaces"
	"bos/types"
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const maxContractInteractions = 50

type ContractInteractionReader struct {
	network interfaces.INetwork
}

func NewContractInteractionReader(network interfaces.INetwork) *ContractInteractionReader {
	return &ContractInteractionReader{
		network: network,
	}
}

func (r *ContractInteractionReader) RecentInteractions(ctx context.Context, contract common.Address, since time.Time) ([]types.ContractInteraction, error) {
	if r.network == nil {
		return nil, fmt.Errorf("network provider is not configured")
	}

	client, _, networkCtx, cancel, err := r.network.Active()
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer client.Close()

	if deadline, ok := ctx.Deadline(); ok {
		var cancelWithDeadline context.CancelFunc
		networkCtx, cancelWithDeadline = context.WithDeadline(networkCtx, deadline)
		defer cancelWithDeadline()
	}

	latest, err := client.HeaderByNumber(networkCtx, nil)
	if err != nil {
		return nil, err
	}

	chainNow := time.Unix(int64(latest.Time), 0)
	if since.After(chainNow) {
		since = chainNow.Add(-24 * time.Hour)
	}

	fromBlock, err := blockAtOrAfter(networkCtx, client, since.Unix(), latest.Number.Uint64())
	if err != nil {
		return nil, err
	}

	return directTransactionInteractions(networkCtx, client, fromBlock, latest.Number.Uint64(), contract, chainNow)
}

func blockAtOrAfter(ctx context.Context, client *ethclient.Client, since int64, latest uint64) (uint64, error) {
	var low uint64
	high := latest
	best := latest

	for low <= high {
		mid := low + (high-low)/2
		header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(mid))
		if err != nil {
			return 0, err
		}

		if int64(header.Time) >= since {
			best = mid
			if mid == 0 {
				break
			}
			high = mid - 1
			continue
		}
		low = mid + 1
	}

	return best, nil
}

func directTransactionInteractions(
	ctx context.Context,
	client *ethclient.Client,
	fromBlock uint64,
	toBlock uint64,
	contract common.Address,
	chainNow time.Time,
) ([]types.ContractInteraction, error) {
	interactions := []types.ContractInteraction{}

	for blockNumber := toBlock; blockNumber >= fromBlock; blockNumber-- {
		block, err := rpcBlockByNumber(ctx, client, blockNumber)
		if err != nil {
			return nil, err
		}

		timestamp, err := hexUint64(block.Timestamp)
		if err != nil {
			return nil, err
		}

		for _, tx := range block.Transactions {
			if !strings.EqualFold(tx.To, contract.Hex()) {
				continue
			}

			interactions = append(interactions, types.ContractInteraction{
				Address: common.HexToAddress(tx.From).Hex(),
				Action:  "Call",
				Age:     formatAge(chainNow.Sub(time.Unix(int64(timestamp), 0))),
				TxHash:  tx.Hash,
			})

			if len(interactions) >= maxContractInteractions {
				return interactions, nil
			}
		}

		if blockNumber == 0 {
			break
		}
	}

	return interactions, nil
}

type rpcBlock struct {
	Timestamp    string              `json:"timestamp"`
	Transactions []rpcTransactionRef `json:"transactions"`
}

type rpcTransactionRef struct {
	Hash string `json:"hash"`
	From string `json:"from"`
	To   string `json:"to"`
}

func rpcBlockByNumber(ctx context.Context, client *ethclient.Client, blockNumber uint64) (rpcBlock, error) {
	var block rpcBlock
	err := client.Client().CallContext(ctx, &block, "eth_getBlockByNumber", fmt.Sprintf("0x%x", blockNumber), true)
	return block, err
}

func hexUint64(value string) (uint64, error) {
	return strconv.ParseUint(strings.TrimPrefix(value, "0x"), 16, 64)
}

func formatAge(age time.Duration) string {
	if age < time.Minute {
		return "<1m"
	}
	if age < time.Hour {
		return fmt.Sprintf("%dm", int(age.Minutes()))
	}
	return fmt.Sprintf("%dh", int(age.Hours()))
}
