package interfaces

import "github.com/ethereum/go-ethereum/accounts"

type IWallet interface {
	Account() (*accounts.Account, error)
	Open() (*accounts.Wallet, *accounts.Account, error)
	SendTransaction(to string, amount *string, gasLimit *uint64) (*string, error)
}
