package interfaces

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type IWallet interface {
	Account() (*common.Address, error)
	Open() (*accounts.Wallet, *accounts.Account, error)
	SendTransaction(to string, amount *string, gasLimit *uint64) (*string, error)
	SignTransaction(ctx context.Context, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}
