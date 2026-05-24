package di

import (
	"bos/implementations"
	"bos/interfaces"
	"bos/repositories"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

var wallet interfaces.IWallet
var network interfaces.INetwork
var logger interfaces.ILogger
var swaps []interfaces.ISwapRouter
var storage interfaces.IStorage
var register repositories.RepositoryRegister

type Dependencies struct {
	Wallet   interfaces.IWallet
	Network  interfaces.INetwork
	Logger   interfaces.ILogger
	Swaps    []interfaces.ISwapRouter
	Storage  interfaces.IStorage
	Register *repositories.RepositoryRegister
}

func SetupDependencies() {
	setupDatabase()
	register = repositories.NewRegister(storage)
	network = implementations.NewNetworkProvider(register.Network)
	wallet = implementations.NewLedger(network, register.Accounts)
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

func SetupDependenciesWith(deps Dependencies) {
	wallet = deps.Wallet
	network = deps.Network
	logger = deps.Logger
	swaps = deps.Swaps
	storage = deps.Storage

	if deps.Register != nil {
		register = *deps.Register
		return
	}

	register = repositories.RepositoryRegister{}
}

func setupDatabase() {
	log.Println("Initializing database")
	storage = implementations.NewStorage(
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

func Repositories() repositories.RepositoryRegister {
	return register
}
