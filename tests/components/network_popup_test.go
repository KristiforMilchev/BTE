package components_test

import (
	"context"
	"math/big"
	"strings"
	"testing"

	networksPopup "bos/components/network_popup"
	"bos/di"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type popupNetwork struct {
	active   types.Network
	networks []types.Network
}

func (m *popupNetwork) Set(_ *string, _ *string, _ *string, _ big.Int) error {
	return nil
}

func (m *popupNetwork) Networks() (*[]types.Network, error) {
	return &m.networks, nil
}

func (m *popupNetwork) Change(network *types.Network) {
	m.active = *network
}

func (m *popupNetwork) Active() (*ethclient.Client, *big.Int, context.Context, context.CancelFunc, error) {
	return nil, big.NewInt(1337), context.Background(), func() {}, nil
}

func (m *popupNetwork) Balance(address common.Address) (*types.NetworkBalanace, error) {
	return &types.NetworkBalanace{
		Address: address.Hex(),
		Balance: "0",
		ChainID: big.NewInt(1337),
	}, nil
}

func (m *popupNetwork) Network() types.Network {
	return m.active
}

func TestNetworkPopupSearchFiltersAndSelectionSubmitsNetwork(t *testing.T) {
	ethereumName := "Ethereum"
	ethereumRPC := "https://eth.llamarpc.com"
	baseName := "Base"
	baseRPC := "https://mainnet.base.org"
	symbol := "ETH"

	network := &popupNetwork{
		networks: []types.Network{
			{Name: &ethereumName, Rpc: &ethereumRPC, Symbol: &symbol, Chain: big.NewInt(1)},
			{Name: &baseName, Rpc: &baseRPC, Symbol: &symbol, Chain: big.NewInt(8453)},
		},
	}
	network.active = network.networks[0]

	di.SetupDependenciesWith(di.Dependencies{Network: network})
	t.Cleanup(func() {
		di.SetupDependenciesWith(di.Dependencies{})
	})

	model := networksPopup.New()
	model.Visible = true

	for _, input := range []tea.KeyMsg{key("b"), key("a"), key("s"), key("e")} {
		var cmd tea.Cmd
		model, cmd = model.Update(input)
		if cmd != nil {
			_ = cmd()
		}
	}

	liveView := model.View()
	if !strings.Contains(liveView, "Base") {
		t.Fatalf("live filtered view = %q, want Base network", liveView)
	}
	if strings.Contains(liveView, "Ethereum") {
		t.Fatalf("live filtered view = %q, want Ethereum filtered out", liveView)
	}

	var cmd tea.Cmd
	model, cmd = model.Update(key("enter"))
	if cmd == nil {
		t.Fatal("search enter returned nil command, want search submission")
	}
	msg := cmd()

	model, cmd = model.Update(msg)
	if cmd != nil {
		_ = cmd()
	}

	view := model.View()
	if !strings.Contains(view, "Base") {
		t.Fatalf("filtered view = %q, want Base network", view)
	}
	if strings.Contains(view, "Ethereum") {
		t.Fatalf("filtered view = %q, want Ethereum filtered out", view)
	}

	model, cmd = model.Update(key("space"))
	if cmd == nil {
		t.Fatal("space returned nil command, want network submission")
	}
	submitted, ok := cmd().(networksPopup.SubmittedMsg)
	if !ok {
		t.Fatalf("network submit returned %T, want networksPopup.SubmittedMsg", submitted)
	}

	baseId := big.NewInt(8453)
	if *submitted.Network.Name != "Base" || submitted.Network.Chain.Cmp(baseId) != 0 {
		t.Fatalf("submitted network = %+v, want Base chain 8453", submitted.Network)
	}
}
