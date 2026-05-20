package di

import (
	"bos/implementations"
	"bos/interfaces"
)

var wallet interfaces.IWallet
var network interfaces.INetwork

func SetupDependencies() {
	network = implementations.NewNetworkProvider()
	wallet = implementations.NewLedger(network)
}

func GetWallet() interfaces.IWallet {
	return wallet
}

func GetNetwork() interfaces.INetwork {
	return network
}
