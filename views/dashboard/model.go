package dashboard

import (
	"bos/components/amount"
	"bos/components/contacts"
	networkDialog "bos/components/network_dialog"
	networksPopup "bos/components/network_popup"
	tokenlist "bos/components/token_list"
	transactionPreview "bos/components/transaction_preview"
	transactionsComponent "bos/components/transactions"
	"bos/enums"
	"bos/interfaces"
	"bos/types"
	"bos/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Wallet        interfaces.IWallet
	Address       string
	Balance       string
	ChainID       string
	statusMessage string
}

type Model struct {
	wallet  interfaces.IWallet
	address string

	width  int
	height int

	statusMessage string

	focus enums.FocusArea

	contacts     *contacts.Model
	transaction  *transactionPreview.Model
	amount       *amount.Model
	tokenList    *tokenlist.Model
	transactions *transactionsComponent.Model

	networkDialog *networkDialog.Model
	networkPopup  *networksPopup.Model
}

func New(config Config) *Model {

	model := &Model{
		wallet:        config.Wallet,
		address:       config.Address,
		focus:         enums.FocusSend,
		contacts:      contacts.NewContacts(),
		amount:        amount.New(),
		tokenList:     tokenlist.New(),
		transactions:  transactionsComponent.New(),
		transaction:   transactionPreview.New(6),
		networkDialog: networkDialog.New(),
		networkPopup:  networksPopup.New(),
	}
	model.transactions.Load()
	return model
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) onNetworkChanged() {
	m.tokenList = tokenlist.New()

}

func (m *Model) OnTransactionSent() {
	m.amount.Clear()
	m.transactions.Load()
	m.statusMessage = "Transaction sent"
	m.focus = enums.FocusSend
}

func (m *Model) beginSend() (tea.Model, tea.Cmd) {
	draft, ok := m.buildDraft()
	if !ok {
		return m, nil
	}

	return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenConfirm, Payload: draft} }
}

func (m *Model) beginSimulation() (tea.Model, tea.Cmd) {
	draft, ok := m.buildDraft()
	if !ok {
		return m, nil
	}

	m.statusMessage = "Simulation started"
	return m, func() tea.Msg {
		return types.NavigateMsg{
			Screen: enums.ScreenSimulationReport,
			Payload: types.SimulationReportPayload{
				Draft:  &draft,
				Return: "dashboard",
			},
		}
	}
}

func (m *Model) buildDraft() (types.TxDraft, bool) {
	amount := strings.TrimSpace(m.amount.Value())
	if amount == "" {
		m.statusMessage = "Enter an amount before sending"
		m.focus = enums.FocusAmount
		return types.TxDraft{}, false
	}

	if _, err := utils.ParseEtherToWei(amount); err != nil {
		m.statusMessage = "Invalid amount: " + err.Error()
		m.focus = enums.FocusAmount
		return types.TxDraft{}, false
	}
	if !m.tokenList.SelectedAsset().Native {
		m.statusMessage = "Token transfer signing is not integrated yet"
		return types.TxDraft{}, false
	}
	if !common.IsHexAddress(m.contacts.SelectedRecipient().Address) {
		m.statusMessage = "Selected contact has an invalid address"
		m.focus = enums.FocusContacts
		return types.TxDraft{}, false
	}

	if !common.IsHexAddress(m.address) {
		m.statusMessage = "Wallet address is not available"
		return types.TxDraft{}, false
	}

	return types.TxDraft{
		FromAddress: m.address, RecipientName: m.contacts.SelectedRecipient().Name, RecipientAddress: m.contacts.SelectedRecipient().Address,
		Amount: amount, Asset: m.tokenList.SelectedAsset(), EstimatedFee: m.transaction.EstimatedFee(),
		SimulationStatus: m.transaction.SimulationStatus(), RiskLevel: m.transaction.RiskLevel(),
	}, true
}
