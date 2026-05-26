package types

import (
	"math/big"

	"github.com/google/uuid"
)

type Network struct {
	Id       uuid.UUID
	Name     *string
	Symbol   *string
	Rpc      *string
	Chain    *big.Int
	Explorer *string
}
