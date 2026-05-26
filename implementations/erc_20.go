package implementations

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const erc20ABI = `[
  {"name":"name","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"string"}]},
  {"name":"symbol","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"string"}]},
  {"name":"decimals","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"uint8"}]},
  {"name":"balanceOf","type":"function","stateMutability":"view","inputs":[{"name":"owner","type":"address"}],"outputs":[{"type":"uint256"}]},
  {"name":"allowance","type":"function","stateMutability":"view","inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"outputs":[{"type":"uint256"}]},
  {"name":"approve","type":"function","stateMutability":"nonpayable","inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}]},
  {"name":"transfer","type":"function","stateMutability":"nonpayable","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}]}
]`

type ERC20 struct {
	address common.Address
	client  *ContractClient
	abi     abi.ABI
}

func NewERC20(client *ContractClient, address common.Address) (*ERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return nil, err
	}

	return &ERC20{
		address: address,
		client:  client,
		abi:     parsed,
	}, nil
}

func (e *ERC20) Address() common.Address {
	return e.address
}

func (e *ERC20) Name(ctx context.Context) (string, error) {
	out, err := e.client.Call(ctx, e.address, e.abi, "name")
	if err != nil {
		return "", err
	}
	return out[0].(string), nil
}

func (e *ERC20) Symbol(ctx context.Context) (string, error) {
	out, err := e.client.Call(ctx, e.address, e.abi, "symbol")
	if err != nil {
		return "", err
	}
	return out[0].(string), nil
}

func (e *ERC20) Decimals(ctx context.Context) (uint8, error) {
	out, err := e.client.Call(ctx, e.address, e.abi, "decimals")
	if err != nil {
		return 0, err
	}
	return out[0].(uint8), nil
}

func (e *ERC20) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	out, err := e.client.Call(ctx, e.address, e.abi, "balanceOf", owner)
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

func (e *ERC20) Allowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
	out, err := e.client.Call(ctx, e.address, e.abi, "allowance", owner, spender)
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

func (e *ERC20) Approve(ctx context.Context, spender common.Address, amount *big.Int) (common.Hash, error) {
	return e.client.Transact(ctx, e.address, e.abi, "approve", nil, spender, amount)
}

func (e *ERC20) Transfer(ctx context.Context, to common.Address, amount *big.Int) (common.Hash, error) {
	return e.client.Transact(ctx, e.address, e.abi, "transfer", nil, to, amount)
}
