package implementations

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const uniswapV2RouterABI = `[
  {"name":"getAmountsOut","type":"function","stateMutability":"view","inputs":[{"name":"amountIn","type":"uint256"},{"name":"path","type":"address[]"}],"outputs":[{"name":"amounts","type":"uint256[]"}]},
  {"name":"swapExactTokensForTokens","type":"function","stateMutability":"nonpayable","inputs":[{"name":"amountIn","type":"uint256"},{"name":"amountOutMin","type":"uint256"},{"name":"path","type":"address[]"},{"name":"to","type":"address"},{"name":"deadline","type":"uint256"}],"outputs":[{"name":"amounts","type":"uint256[]"}]}
]`

type UniswapV2Router struct {
	address common.Address
	client  *ContractClient
	abi     abi.ABI
}

func NewUniswapV2Router(client *ContractClient, address common.Address) (*UniswapV2Router, error) {
	parsed, err := abi.JSON(strings.NewReader(uniswapV2RouterABI))
	if err != nil {
		return nil, err
	}

	return &UniswapV2Router{
		address: address,
		client:  client,
		abi:     parsed,
	}, nil
}

func (r *UniswapV2Router) Address() common.Address {
	return r.address
}

func (r *UniswapV2Router) GetAmountsOut(
	ctx context.Context,
	amountIn *big.Int,
	path []common.Address,
) ([]*big.Int, error) {
	out, err := r.client.Call(ctx, r.address, r.abi, "getAmountsOut", amountIn, path)
	if err != nil {
		return nil, err
	}

	return out[0].([]*big.Int), nil
}

func (r *UniswapV2Router) SwapExactTokensForTokens(
	ctx context.Context,
	amountIn *big.Int,
	amountOutMin *big.Int,
	path []common.Address,
	to common.Address,
	deadline *big.Int,
) (common.Hash, error) {
	return r.client.Transact(
		ctx,
		r.address,
		r.abi,
		"swapExactTokensForTokens",
		nil,
		amountIn,
		amountOutMin,
		path,
		to,
		deadline,
	)
}
