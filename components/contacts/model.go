package contacts

import (
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	contacts        []types.Contact
	selectedContact int
}

func (m Model) SelectedRecipient() types.Contact {
	if len(m.contacts) == 0 {
		return types.Contact{Name: "No Contact", Address: ""}
	}
	if m.selectedContact < 0 || m.selectedContact >= len(m.contacts) {
		return m.contacts[0]
	}
	return m.contacts[m.selectedContact]
}

func (m *Model) Init() tea.Cmd { return nil }

func NewContacts() *Model {
	contacts := []types.Contact{
		{Name: "Treasury Wallet", Address: "0x1111111111111111111111111111111111111111"},
		{Name: "Personal Wallet", Address: "0x2222222222222222222222222222222222222222"},
		{Name: "Binance Deposit", Address: "0x3333333333333333333333333333333333333333"},
	}
	return &Model{
		contacts:        contacts,
		selectedContact: 0,
	}
}
