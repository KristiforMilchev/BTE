package dashboard

import (
	"bos/enums"
	"bos/types"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
		return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenLoading} }

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
	case "up", "k", "down", "j":
		if m.focus == enums.FocusTokens {
			_, cmd := m.tokenList.Update(msg)
			m.resetAnalysis()
			return m, cmd
		}
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

func (m *Model) activateFocusedItem() (tea.Model, tea.Cmd) {
	switch m.focus {
	case enums.FocusContacts:
		m.focus = enums.FocusAmount
		m.syncFocus()
		m.statusMessage = "Recipient selected: " + m.contacts.SelectedRecipient().Name
	case enums.FocusTokens:
		m.focus = enums.FocusAmount
		m.syncFocus()
		m.statusMessage = "Asset selected: " + m.tokenList.SelectedAsset().Symbol
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
