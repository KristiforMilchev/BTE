package activeWalletLabel

import (
	"bos/di"
	"log"
)

type Model struct {
	wallet string
}

func New() *Model {

	address, err := di.GetWallet().Account()
	if err != nil {
		log.Printf("Can't connect to wallet to retrive wallet label")
	}

	return &Model{
		wallet: address.Hex(),
	}
}
