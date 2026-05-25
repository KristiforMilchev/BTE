package networksPopup

import (
	"strings"

	"bos/components/search"

	tea "github.com/charmbracelet/bubbletea"
)

type focusArea int

const (
	focusSearch focusArea = iota
	focusTable
)

type Network struct {
	Name    string
	RPC     string
	Symbol  string
	ChainID int64
}

type SubmittedMsg struct {
	Network Network
}

type CancelledMsg struct{}

type Model struct {
	Visible bool

	search *search.Model
	focus  focusArea

	networks []Network
	filtered []Network
	selected int
	offset   int
}

func (m *Model) Init() tea.Cmd {
	if m.search == nil {
		return nil
	}
	return m.search.Init()
}

func New() *Model {
	networks := defaultNetworks()
	return &Model{
		Visible:  false,
		search:   search.New("Search networks", 64),
		focus:    focusSearch,
		networks: networks,
		filtered: networks,
	}
}

func (m *Model) focusSearch() {
	m.focus = focusSearch
	m.search.Focus()
}

func (m *Model) focusTable() {
	m.focus = focusTable
	m.search.Blur()
}

func (m *Model) applySearch(query string) {
	m.filterNetworks(query)
	m.focusTable()
}

func (m *Model) filterNetworks(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		m.filtered = m.networks
	} else {
		filtered := make([]Network, 0, len(m.networks))
		for _, network := range m.networks {
			if strings.Contains(strings.ToLower(network.Name), query) ||
				strings.Contains(strings.ToLower(network.RPC), query) ||
				strings.Contains(strings.ToLower(network.Symbol), query) ||
				strings.Contains(strings.ToLower(chainIDString(network.ChainID)), query) {
				filtered = append(filtered, network)
			}
		}
		m.filtered = filtered
	}

	m.selected = 0
	m.offset = 0
}

func (m *Model) submit() tea.Cmd {
	return func() tea.Msg {
		if len(m.filtered) == 0 {
			return nil
		}
		return SubmittedMsg{Network: m.filtered[m.selected]}
	}
}

func defaultNetworks() []Network {
	return []Network{
		{Name: "Blockcert", RPC: "https://rpc.blockcert.net", Symbol: "ETH", ChainID: 707},
		{Name: "Ethereum", RPC: "https://eth.llamarpc.com", Symbol: "ETH", ChainID: 1},
		{Name: "Sepolia", RPC: "https://ethereum-sepolia-rpc.publicnode.com", Symbol: "ETH", ChainID: 11155111},
		{Name: "Polygon", RPC: "https://polygon-rpc.com", Symbol: "POL", ChainID: 137},
		{Name: "Arbitrum One", RPC: "https://arb1.arbitrum.io/rpc", Symbol: "ETH", ChainID: 42161},
		{Name: "Optimism", RPC: "https://mainnet.optimism.io", Symbol: "ETH", ChainID: 10},
		{Name: "Base", RPC: "https://mainnet.base.org", Symbol: "ETH", ChainID: 8453},
		{Name: "BSC", RPC: "https://bsc-dataseed.binance.org", Symbol: "BNB", ChainID: 56},
		{Name: "Avalanche", RPC: "https://api.avax.network/ext/bc/C/rpc", Symbol: "AVAX", ChainID: 43114},
		{Name: "Localhost", RPC: "http://127.0.0.1:8545", Symbol: "ETH", ChainID: 31337},
	}
}
