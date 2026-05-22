package dashboard

import (
	"bos/enums"
	"bos/utils"
	"bos/views"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/common"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
		return m, func() tea.Msg { return views.NavigateMsg{Screen: views.ScreenLoading} }
	case "s":
		m.runFakeSimulation()
		return m, nil
	case "S", "shift+s":
		return m.beginSend()
	case "p", "P":
		m.focus = enums.FocusContacts
		m.syncFocus()
		m.statusMessage = "Recipient picker active"
		return m, nil
	case "esc":
		m.focus = enums.FocusAmount
		m.syncFocus()
		m.statusMessage = "Amount editor active"
		return m, nil
	}

	switch msg.String() {
	case "left", "h":
		m.focus = enums.FocusAmount
		m.syncFocus()
		return m, nil
	case "right", "l":
		m.focus = enums.FocusTokens
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

	if m.focus == enums.FocusAmount {
		var cmd tea.Cmd
		m.amountInput, cmd = m.amountInput.Update(msg)
		m.resetAnalysis()
		return m, cmd
	}

	return m, nil
}

func (m *Model) moveSelection(delta int) {
	switch m.focus {
	case enums.FocusTokens:
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
	case enums.FocusContacts:
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
	case enums.FocusAmount:
		if delta > 0 {
			m.focus = enums.FocusContacts
		} else {
			m.focus = enums.FocusTokens
		}
		m.syncFocus()
	}
}

func (m *Model) activateFocusedItem() (tea.Model, tea.Cmd) {
	switch m.focus {
	case enums.FocusContacts:
		m.focus = enums.FocusAmount
		m.syncFocus()
		m.statusMessage = "Recipient selected: " + m.selectedRecipient().Name
	case enums.FocusTokens:
		m.focus = enums.FocusAmount
		m.syncFocus()
		m.statusMessage = "Asset selected: " + m.selectedAsset().Symbol
		m.resetAnalysis()
	case enums.FocusAmount:
		m.statusMessage = "Amount set: " + strings.TrimSpace(m.amountInput.Value())
	}
	return m, nil
}

func (m *Model) syncFocus() {
	if m.focus == enums.FocusAmount {
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

func (m *Model) beginSend() (tea.Model, tea.Cmd) {
	amount := strings.TrimSpace(m.amountInput.Value())
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
	if !m.selectedAsset().Native {
		m.statusMessage = "Token transfer signing is not integrated yet"
		return m, nil
	}
	if !common.IsHexAddress(m.selectedRecipient().Address) {
		m.statusMessage = "Selected contact has an invalid address"
		m.focus = enums.FocusContacts
		m.syncFocus()
		return m, nil
	}
	draft := views.TxDraft{
		FromAddress: m.address, RecipientName: m.selectedRecipient().Name, RecipientAddress: m.selectedRecipient().Address,
		Amount: amount, Asset: m.selectedAsset(), EstimatedFee: m.estimatedFee,
		SimulationStatus: m.simulationStatus, RiskLevel: m.riskLevel,
	}
	return m, func() tea.Msg { return views.NavigateMsg{Screen: views.ScreenConfirm, Payload: draft} }
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
