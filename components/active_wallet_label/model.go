package activeWalletLabel

import (
	"bos/di"
	"log"
)

type Model struct {
	wallet string
}

func New() *Model {
	wallet := di.GetWallet()
	if wallet == nil {
		return &Model{}
	}

	address, err := wallet.Account()
	if err != nil || address == nil {
		log.Printf("Can't connect to wallet to retrive wallet label")
		return &Model{}
	}

	return &Model{
		wallet: address.Hex(),
	}
}
