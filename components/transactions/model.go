package transactions

import (
	transactionComponent "bos/components/transaction"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	transactions []transactionComponent.Model
	selected     int
	offset       int
}

func New(items []types.Transaction) *Model {
	m := &Model{}
	m.SetTransactions(items)
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
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
