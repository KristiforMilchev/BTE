package di_test

import (
	"context"
	"math/big"
	"testing"

	"bos/di"
	"bos/interfaces"
	"bos/repositories"
	"bos/tests/testmocks"
	"bos/types"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type mockWallet struct{}

func (mockWallet) Account() (*common.Address, error) {
	address := common.HexToAddress("0x1111111111111111111111111111111111111111")
	return &address, nil
}

func (mockWallet) Open() (*accounts.Wallet, *accounts.Account, error) {
	return nil, nil, nil
}

func (mockWallet) SendTransaction(_ string, _ *string, _ *uint64) (*string, error) {
	hash := "0xtx"
	return &hash, nil
}

func (mockWallet) SignTransaction(_ context.Context, tx *coretypes.Transaction, _ *big.Int) (*coretypes.Transaction, error) {
	return tx, nil
}

type mockNetwork struct{}

func (mockNetwork) Set(_ *string, _ *string, _ *string, _ big.Int) error {
	return nil
}

func (mockNetwork) Networks() (*[]types.Network, error) {
	return &[]types.Network{}, nil
}

func (mockNetwork) Change(_ *types.Network) {
}

func (mockNetwork) Active() (*ethclient.Client, *big.Int, context.Context, context.CancelFunc, error) {
	return nil, big.NewInt(1337), context.Background(), func() {}, nil
}

func (mockNetwork) Balance(address common.Address) (*types.NetworkBalanace, error) {
	return &types.NetworkBalanace{
		Address: address.Hex(),
		Balance: "1.00000000",
		ChainID: big.NewInt(1337),
	}, nil
}

func (mockNetwork) Network() types.Network {
	name := "Mocknet"
	rpc := "http://mock"
	symbol := "ETH"
	return types.Network{
		Name:   &name,
		Rpc:    &rpc,
		Symbol: &symbol,
		Chain:  big.NewInt(1337),
	}
}

type mockLogger struct{}

func (mockLogger) Write(p []byte) (int, error) {
	return len(p), nil
}

func (mockLogger) Logs() []string {
	return []string{"mock log"}
}

func (mockLogger) Close() error {
	return nil
}

type mockSwapRouter struct {
	address common.Address
}

func (m mockSwapRouter) Address() common.Address {
	return m.address
}

func (mockSwapRouter) GetAmountsOut(_ context.Context, _ *big.Int, _ []common.Address) ([]*big.Int, error) {
	return []*big.Int{big.NewInt(1)}, nil
}

func (mockSwapRouter) SwapExactTokensForTokens(
	_ context.Context,
	_ *big.Int,
	_ *big.Int,
	_ []common.Address,
	_ common.Address,
	_ *big.Int,
) (common.Hash, error) {
	return common.HexToHash("0x1"), nil
}

func TestSetupDependenciesWithUsesProvidedServices(t *testing.T) {
	storage := testmocks.NewStorage(t)
	register := repositories.NewRegister(storage)
	swap := mockSwapRouter{address: common.HexToAddress("0x2222222222222222222222222222222222222222")}

	di.SetupDependenciesWith(di.Dependencies{
		Wallet:   mockWallet{},
		Network:  mockNetwork{},
		Logger:   mockLogger{},
		Swaps:    []interfaces.ISwapRouter{swap},
		Storage:  storage,
		Register: &register,
	})
	t.Cleanup(func() {
		di.SetupDependenciesWith(di.Dependencies{})
	})

	if di.GetWallet() == nil {
		t.Fatal("GetWallet() = nil, want provided wallet")
	}
	if di.GetNetwork() == nil {
		t.Fatal("GetNetwork() = nil, want provided network")
	}
	if di.GetLogger() == nil {
		t.Fatal("GetLogger() = nil, want provided logger")
	}
	if len(di.GetSwaps()) != 1 || di.GetSwaps()[0].Address() != swap.address {
		t.Fatalf("GetSwaps() = %+v, want provided swap", di.GetSwaps())
	}
}
