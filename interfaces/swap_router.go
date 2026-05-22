package interfaces

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ISwapRouter interface {
	Address() common.Address

	GetAmountsOut(
		ctx context.Context,
		amountIn *big.Int,
		path []common.Address,
	) ([]*big.Int, error)

	SwapExactTokensForTokens(
		ctx context.Context,
		amountIn *big.Int,
		amountOutMin *big.Int,
		path []common.Address,
		to common.Address,
		deadline *big.Int,
	) (common.Hash, error)
}
