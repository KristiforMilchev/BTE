package interfaces

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
)

type IWallet interface {
	Account() (*common.Address, error)
	Open() (*accounts.Wallet, *accounts.Account, error)
	SendTransaction(to string, amount *string, gasLimit *uint64) (*string, error)
}
