package dashboard

import (
	"bos/components/amount"
	"bos/components/contacts"
	tokenlist "bos/components/token_list"
	"bos/enums"
	"bos/interfaces"
	"bos/types"
	"bos/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Wallet  interfaces.IWallet
	Address string
	Balance string
	ChainID string
}

type Model struct {
	wallet interfaces.IWallet

	width  int
	height int

	address string
	balance string
	chainID string

	focus enums.FocusArea

	simulationStatus string
	riskLevel        string
	estimatedFee     string
	statusMessage    string
	tokenList        tokenlist.Model
	contacts         contacts.Model
	amount           *amount.Model
}

func New(config Config) *Model {

	return &Model{
		wallet:           config.Wallet,
		address:          config.Address,
		balance:          config.Balance,
		chainID:          config.ChainID,
		focus:            enums.FocusAmount,
		simulationStatus: "Not Run",
		riskLevel:        "—",
		estimatedFee:     "—",
		statusMessage:    "Wallet loaded",
		tokenList:        *tokenlist.New(),
		contacts:         *contacts.NewContacts(),
		amount:           amount.New(),
	}
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) beginSend() (tea.Model, tea.Cmd) {
	amount := strings.TrimSpace(m.amount.Value())
	if amount == "" {
		m.statusMessage = "Enter an amount before sending"
		m.focus = enums.FocusAmount
		m.syncFocus()
		return m, nil
	}
	if _, err := utils.ParseEtherToWei(amount); err != nil {
		m.statusMessage = "Invalid amount: " + err.Error()
		m.focus = enums.FocusAmount
		m.syncFocus()
		return m, nil
	}
	if !m.tokenList.SelectedAsset().Native {
		m.statusMessage = "Token transfer signing is not integrated yet"
		return m, nil
	}
	if !common.IsHexAddress(m.contacts.SelectedRecipient().Address) {
		m.statusMessage = "Selected contact has an invalid address"
		m.focus = enums.FocusContacts
		m.syncFocus()
		return m, nil
	}
	draft := types.TxDraft{
		FromAddress: m.address, RecipientName: m.contacts.SelectedRecipient().Name, RecipientAddress: m.contacts.SelectedRecipient().Address,
		Amount: amount, Asset: m.tokenList.SelectedAsset(), EstimatedFee: m.estimatedFee,
		SimulationStatus: m.simulationStatus, RiskLevel: m.riskLevel,
	}
	return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenConfirm, Payload: draft} }
}

func (m *Model) runFakeSimulation() {
	amount := strings.TrimSpace(m.amount.Value())
	if amount == "" {
		m.statusMessage = "Enter an amount before simulation"
		m.focus = enums.FocusAmount
		m.syncFocus()
		return
	}
	if _, err := utils.ParseEtherToWei(amount); err != nil {
		m.statusMessage = "Invalid amount: " + err.Error()
		m.focus = enums.FocusAmount
		m.syncFocus()
		return
	}
	m.simulationStatus = "Completed"
	m.riskLevel = "Low"
	m.estimatedFee = "0.000021 ETH"
	m.statusMessage = "Simulation completed with fake analysis data"
}
