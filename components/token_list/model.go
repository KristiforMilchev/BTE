package tokenlist

import (
	"bos/di"
	"bos/types"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tokens        []types.Token
	selectedToken int
	offset        int
}

func (m *Model) Init() tea.Cmd {

	return nil

}

func New() *Model {
	m := &Model{}

	network := di.GetNetwork()
	address, err := di.GetWallet().Account()
	if err != nil {
		log.Printf("Can't find account to pin balance -> %s", err)
	}
	details, err := network.Balance(*address)
	if err != nil {
		log.Printf("Get read network balance for address -> %s | %s", address, err)
	}

	log.Printf("m.balance: %s\n", details.Balance)

	tokens := []types.Token{
		{Symbol: "ETH", Name: details.Address, Balance: details.Balance, Address: "native", Native: true},
	}
	m.tokens = tokens
	log.Println("Updating network balance")
	return m
}

func (m *Model) SetTokens(tokens []types.Token) {
	m.tokens = tokens
	if len(m.tokens) == 0 {
		m.selectedToken = 0
		m.offset = 0
		return
	}
	if m.selectedToken >= len(m.tokens) {
		m.selectedToken = len(m.tokens) - 1
	}
	if m.selectedToken < 0 {
		m.selectedToken = 0
	}
}
