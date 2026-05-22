package main

import (
	"bos/utils"
	"bos/views"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/common"
)

func (m *Model) Init() tea.Cmd {
	return m.loadWallet()
}

func (m *Model) loadWallet() tea.Cmd {
	return func() tea.Msg {
		wallet, err := m.wallet.Account()
		if err != nil {
			return walletLoadedMsg{Err: err}
		}

		networkBalance, err := m.network.Balance(wallet.Address)
		if err != nil {
			return walletLoadedMsg{Err: err}
		}

		return walletLoadedMsg{
			Address: wallet.Address.Hex(),
			Balance: networkBalance.Balance,
			ChainID: networkBalance.ChainID.String(),
		}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case walletLoadedMsg:
		if msg.Err != nil {
			m.screen = views.ScreenError
			m.err = msg.Err.Error()
			return m, nil
		}

		m.address = msg.Address
		m.balance = msg.Balance
		m.chainID = msg.ChainID
		m.tokens[0].Balance = msg.Balance
		m.screen = views.ScreenDashboard
		m.statusMessage = "Wallet loaded"
		return m, nil

	case sendFinishedMsg:
		if msg.Err != nil {
			m.screen = views.ScreenError
			m.err = msg.Err.Error()
			return m, nil
		}

		m.txHash = msg.TxHash
		m.screen = views.ScreenSent
		m.statusMessage = "Transaction broadcast"
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "r":
		m.screen = views.ScreenLoading
		m.err = ""
		m.txHash = ""
		m.statusMessage = "Refreshing wallet"
		return m, m.loadWallet()
	}

	switch m.screen {
	case views.ScreenDashboard:
		return m.handleDashboardKey(msg)
	case views.ScreenConfirm:
		return m.handleConfirmKey(msg)
	case views.ScreenSent:
		if msg.String() == "enter" || msg.String() == "esc" {
			m.screen = views.ScreenDashboard
			m.txHash = ""
			return m, nil
		}
	case views.ScreenError:
		if msg.String() == "enter" {
			m.screen = views.ScreenLoading
			m.err = ""
			return m, m.loadWallet()
		}
		if msg.String() == "esc" {
			m.screen = views.ScreenDashboard
			m.err = ""
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) handleDashboardKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Direct action keys. Keep sending/simulation out of the navigation model.
	switch key {
	case "s":
		m.runFakeSimulation()
		return m, nil
	case "S", "shift+s":
		return m.beginSend()
	case "p", "P":
		m.focus = views.FocusContacts
		m.syncFocus()
		m.statusMessage = "Recipient picker active"
		return m, nil
	case "esc":
		m.focus = views.FocusAmount
		m.syncFocus()
		m.statusMessage = "Amount editor active"
		return m, nil
	}

	// Spatial navigation.
	switch key {
	case "left", "h":
		m.focus = views.FocusAmount
		m.syncFocus()
		return m, nil
	case "right", "l":
		m.focus = views.FocusTokens
		m.syncFocus()
		return m, nil
	case "up", "k":
		m.moveSelection(-1)
		return m, nil
	case "down", "j":
		m.moveSelection(1)
		return m, nil
	case "enter", " ":
		return m.activateFocusedItem()
	}

	if m.focus == views.FocusAmount {
		var cmd tea.Cmd
		m.amountInput, cmd = m.amountInput.Update(msg)
		m.resetAnalysis()
		return m, cmd
	}

	return m, nil
}

func (m *Model) moveSelection(delta int) {
	switch m.focus {
	case views.FocusTokens:
		if len(m.tokens) == 0 {
			return
		}
		m.selectedToken += delta
		if m.selectedToken < 0 {
			m.selectedToken = len(m.tokens) - 1
		}
		if m.selectedToken >= len(m.tokens) {
			m.selectedToken = 0
		}
		m.resetAnalysis()
		m.statusMessage = "Selected asset: " + m.selectedAsset().Symbol
	case views.FocusContacts:
		if len(m.contacts) == 0 {
			return
		}
		m.selectedContact += delta
		if m.selectedContact < 0 {
			m.selectedContact = len(m.contacts) - 1
		}
		if m.selectedContact >= len(m.contacts) {
			m.selectedContact = 0
		}
		m.resetAnalysis()
		m.statusMessage = "Selected recipient: " + m.selectedRecipient().Name
	case views.FocusAmount:
		if delta > 0 {
			m.focus = views.FocusContacts
		} else {
			m.focus = views.FocusTokens
		}
		m.syncFocus()
	}
}

func (m *Model) activateFocusedItem() (tea.Model, tea.Cmd) {
	switch m.focus {
	case views.FocusContacts:
		m.focus = views.FocusAmount
		m.syncFocus()
		m.statusMessage = "Recipient selected: " + m.selectedRecipient().Name
		return m, nil
	case views.FocusTokens:
		m.focus = views.FocusAmount
		m.syncFocus()
		m.statusMessage = "Asset selected: " + m.selectedAsset().Symbol
		m.resetAnalysis()
		return m, nil
	case views.FocusAmount:
		m.statusMessage = "Amount set: " + strings.TrimSpace(m.amountInput.Value())
		return m, nil
	default:
		return m, nil
	}
}

func (m *Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter", " ":
		m.screen = views.ScreenSending
		m.statusMessage = "Waiting for Ledger approval"
		return m, m.sendTransaction()
	case "n", "esc":
		m.screen = views.ScreenDashboard
		m.statusMessage = "Transaction cancelled"
		return m, nil
	}

	return m, nil
}

func (m *Model) syncFocus() {
	if m.focus == views.FocusAmount {
		m.amountInput.Focus()
		return
	}
	m.amountInput.Blur()
}

func (m *Model) resetAnalysis() {
	m.simulationStatus = "Not Run"
	m.riskLevel = "—"
	m.estimatedFee = "—"
}

func (m *Model) runFakeSimulation() {
	amount := strings.TrimSpace(m.amountInput.Value())
	if amount == "" {
		m.statusMessage = "Enter an amount before simulation"
		m.focus = views.FocusAmount
		m.syncFocus()
		return
	}
	if _, err := utils.ParseEtherToWei(amount); err != nil {
		m.statusMessage = "Invalid amount: " + err.Error()
		m.focus = views.FocusAmount
		m.syncFocus()
		return
	}

	m.simulationStatus = "Completed"
	m.riskLevel = "Low"
	m.estimatedFee = "0.000021 ETH"
	m.statusMessage = "Simulation completed with fake analysis data"
}

func (m *Model) beginSend() (tea.Model, tea.Cmd) {
	amount := strings.TrimSpace(m.amountInput.Value())
	if amount == "" {
		m.statusMessage = "Enter an amount before sending"
		m.focus = views.FocusAmount
		m.syncFocus()
		return m, nil
	}

	if _, err := utils.ParseEtherToWei(amount); err != nil {
		m.statusMessage = "Invalid amount: " + err.Error()
		m.focus = views.FocusAmount
		m.syncFocus()
		return m, nil
	}

	if !m.selectedAsset().Native {
		m.statusMessage = "Token transfer signing is not integrated yet"
		return m, nil
	}

	if !common.IsHexAddress(m.selectedRecipient().Address) {
		m.statusMessage = "Selected contact has an invalid address"
		m.focus = views.FocusContacts
		m.syncFocus()
		return m, nil
	}

	m.screen = views.ScreenConfirm
	return m, nil
}

func (m *Model) sendTransaction() tea.Cmd {
	recipient := m.selectedRecipient().Address
	amount := strings.TrimSpace(m.amountInput.Value())

	return func() tea.Msg {
		txHash, err := m.wallet.SendTransaction(recipient, &amount, nil)
		if err != nil {
			return sendFinishedMsg{Err: err}
		}
		if txHash == nil {
			return sendFinishedMsg{Err: fmt.Errorf("wallet returned an empty transaction hash")}
		}
		return sendFinishedMsg{TxHash: *txHash}
	}
}

func (m Model) selectedAsset() views.Token {
	if len(m.tokens) == 0 {
		return views.Token{Symbol: "ETH", Balance: "0", Native: true}
	}
	if m.selectedToken < 0 || m.selectedToken >= len(m.tokens) {
		return m.tokens[0]
	}
	return m.tokens[m.selectedToken]
}

func (m Model) selectedRecipient() views.Contact {
	if len(m.contacts) == 0 {
		return views.Contact{Name: "No Contact", Address: ""}
	}
	if m.selectedContact < 0 || m.selectedContact >= len(m.contacts) {
		return m.contacts[0]
	}
	return m.contacts[m.selectedContact]
}
