package implementations

import (
	"bos/constants"
	"bos/repositories"
	"bos/types"
	"bos/utils"
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Network struct {
	network           types.Network
	networkRepository repositories.NetworkRepository
}

func (n *Network) Set(rpc *string, name *string, symbol *string, chain big.Int) error {
	n.network = types.Network{
		Name:   name,
		Symbol: symbol,
		Rpc:    rpc,
		Chain:  &chain,
	}
	return nil
}

func (n *Network) Networks() (*[]types.Network, error) {
	return n.networkRepository.Networks()
}

func (n *Network) Change(network *types.Network) {
	n.network = *network
}

func (n *Network) Active() (*ethclient.Client, *big.Int, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	if constants.RpcURL == "" {
		cancel()
		return nil, nil, nil, nil, errors.New("RPC URL is empty")
	}

	client, err := ethclient.DialContext(ctx, constants.RpcURL)
	if err != nil {
		cancel()
		return nil, nil, nil, nil, fmt.Errorf("failed to connect to RPC %q: %w", constants.RpcURL, err)
	}

	if client == nil {
		cancel()
		return nil, nil, nil, nil, fmt.Errorf("failed to connect to RPC %q: ethclient returned nil client", constants.RpcURL)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		client.Close()
		cancel()
		return nil, nil, nil, nil, fmt.Errorf("failed to read chain ID from RPC %q: %w", constants.RpcURL, err)
	}

	return client, chainID, ctx, cancel, nil
}

func (n *Network) Balance(address common.Address) (*types.NetworkBalanace, error) {
	client, chain, ctx, cancel, err := n.Active()
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer client.Close()

	balance, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	ether := utils.WeiToEther(balance)

	return &types.NetworkBalanace{
		Address: address.Hex(),
		Balance: ether,
		ChainID: chain,
	}, nil
}

func (n *Network) Network() types.Network {
	return n.network
}

func NewNetworkProvider(networkRepository repositories.NetworkRepository) *Network {
	rpc := constants.RpcURL
	symbol := "BGC"
	name := "Blockcert"
	return &Network{
		network: types.Network{
			Name:   &name,
			Symbol: &symbol,
			Rpc:    &rpc,
			Chain:  big.NewInt(31337),
		},
		networkRepository: networkRepository,
	}
}
