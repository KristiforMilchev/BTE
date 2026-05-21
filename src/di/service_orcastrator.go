package di

import (
	"bos/implementations"
	"bos/interfaces"
)

var wallet interfaces.IWallet
var network interfaces.INetwork
var logger interfaces.ILogger

func SetupDependencies() {
	network = implementations.NewNetworkProvider()
	wallet = implementations.NewLedger(network)
	loggerInstance, err := implementations.New("blockcert.log")
	if err != nil {
		panic(err)
	}

	logger = loggerInstance
}

func GetWallet() interfaces.IWallet {
	return wallet
}

func GetNetwork() interfaces.INetwork {
	return network
}
func GetLogger() interfaces.ILogger {
	return logger
}
