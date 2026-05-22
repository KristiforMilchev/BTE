package tokenlist

import (
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tokens        []types.Token
	selectedToken int
}

func (m *Model) Init() tea.Cmd { return nil }

func New() *Model {
	tokens := []types.Token{
		{Symbol: "ETH", Name: "Native Network Token", Balance: "0", Address: "native", Native: true},
		{Symbol: "USDT", Name: "Tether USD", Balance: "1240.22", Address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"},
		{Symbol: "LINK", Name: "Chainlink", Balance: "4.12", Address: "0x514910771af9ca656af840dff83e8264ecf986ca"},
		{Symbol: "UNI", Name: "Uniswap", Balance: "42.11", Address: "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},
	}

	return &Model{
		tokens: tokens,
	}
}
