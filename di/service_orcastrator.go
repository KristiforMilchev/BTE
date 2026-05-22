package di

import (
	"bos/implementations"
	"bos/interfaces"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

var wallet interfaces.IWallet
var network interfaces.INetwork
var logger interfaces.ILogger
var swaps []interfaces.ISwapRouter

func SetupDependencies() {
	setupDatabase()
	network = implementations.NewNetworkProvider()
	wallet = implementations.NewLedger(network)
	swapRouter := implementations.NewContractClient(network, nil)
	router := common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")
	swap, err := implementations.NewUniswapV2Router(swapRouter, router)
	if err != nil {
		log.Printf("Failed initialiaze a swap -> %s", router)
		return
	}

	swaps = []interfaces.ISwapRouter{}
	swaps = append(swaps, swap)

	loggerInstance, err := implementations.New("blockcert.log")
	if err != nil {
		panic(err)
	}

	logger = loggerInstance
}

func setupDatabase() {
	log.Println("Initializing database")
	storage := implementations.NewStorage(
		"./data/bos.db",
		"./sql/schema.sql",
	)

	if err := storage.Init(); err != nil {
		panic(err)
	}
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

func GetSwaps() []interfaces.ISwapRouter {
	return swaps
}
