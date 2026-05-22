package interfaces

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type IERC20 interface {
	Address() common.Address

	Name(ctx context.Context) (string, error)
	Symbol(ctx context.Context) (string, error)
	Decimals(ctx context.Context) (uint8, error)

	BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error)
	Allowance(ctx context.Context, owner, spender common.Address) (*big.Int, error)

	Approve(ctx context.Context, spender common.Address, amount *big.Int) (common.Hash, error)
	Transfer(ctx context.Context, to common.Address, amount *big.Int) (common.Hash, error)
	TransferFrom(ctx context.Context, from, to common.Address, amount *big.Int) (common.Hash, error)
}
