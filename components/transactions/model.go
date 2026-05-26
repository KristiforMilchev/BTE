package transactions

import (
	transactionComponent "bos/components/transaction"
	"bos/di"
	"bos/types"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	transactions []transactionComponent.Model
	selected     int
	offset       int
}

func New(items ...[]types.Transaction) *Model {
	m := &Model{}
	if len(items) > 0 {
		m.SetTransactions(items[0])
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Load() {
	transactionsRepository := di.Repositories().Transactions
	wallet, err := di.GetWallet().Account()
	if err != nil {
		log.Printf("can't get transactions, account not set -> %s", err)
		return
	}

	network := di.GetNetwork().Network()
	walletHex := wallet.Hex()
	items, err := transactionsRepository.GetTransactions(&network.Id, &walletHex)
	if err != nil {
		return
	}
	m.SetTransactions(*items)
}

func (m *Model) SetTransactions(items []types.Transaction) {
	m.transactions = make([]transactionComponent.Model, 0, len(items))
	for _, item := range items {
		m.transactions = append(m.transactions, transactionComponent.New(item))
	}
	if len(m.transactions) == 0 {
		m.selected = 0
		m.offset = 0
		return
	}

	if m.selected >= len(m.transactions) {
		m.selected = len(m.transactions) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}
