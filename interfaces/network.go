package interfaces

import (
	"bos/types"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type INetwork interface {
	Set(rpc *string, name *string, symbol *string, chain big.Int) error
	Networks() (*[]types.Network, error)
	Change(rpc *string) error
	Active() (*ethclient.Client, *big.Int, context.Context, context.CancelFunc, error)
	Balance(address common.Address) (*types.NetworkBalanace, error)
	Name() *string
}
