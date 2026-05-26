package interfaces

import (
	"bos/types"
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type IContractInteractionReader interface {
	RecentInteractions(ctx context.Context, contract common.Address, since time.Time) ([]types.ContractInteraction, error)
}
