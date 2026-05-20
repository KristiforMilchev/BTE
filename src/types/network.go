package types

import "math/big"

type Network struct {
	Name   *string
	Symbol *string
	Rpc    *string
	Chain  *big.Int
}
