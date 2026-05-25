package dashboard

import (
	networkDialog "bos/components/network_dialog"
	networksPopup "bos/components/network_popup"
	"bos/enums"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case networkDialog.SubmittedMsg:
		m.networkDialog.Visible = false
		// save msg.Network here
		return m, nil

	case networkDialog.CancelledMsg:
		m.networkDialog.Visible = false
		return m, nil
	case networksPopup.SubmittedMsg:
		m.networkPopup.Visible = false
		return m, nil
	case networksPopup.CancelledMsg:
		m.networkPopup.Visible = false
		return m, nil
	}

	if m.networkDialog.Visible {
		var cmd tea.Cmd
		m.networkDialog, cmd = m.networkDialog.Update(msg)
		return m, cmd
	}

	if m.networkPopup.Visible {
		var cmd tea.Cmd
		m.networkPopup, cmd = m.networkPopup.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.focus == enums.FocusContacts {
		msg, _ := m.contacts.Update(msg)
		if msg != nil {
			m.focus = enums.FocusSend
		}
		return m, nil
	}

	if m.focus == enums.FocusTokens {
		if msg.String() == "l" || msg.String() == "right" {
			m.focus = enums.FocusTransactions
			return m, nil
		}
		msg, _ := m.tokenList.Update(msg)
		if msg != nil {
			m.focus = enums.FocusSend
		}
		return m, nil
	}

	if m.focus == enums.FocusTransactions {
		if msg.String() == "h" || msg.String() == "left" {
			m.focus = enums.FocusTokens
			return m, nil
		}
		msg, _ := m.transactions.Update(msg)
		if msg != nil {
			m.focus = enums.FocusSend
		}
		return m, nil
	}

	if m.focus == enums.FocusAmount {
		msg, _ := m.amount.Update(msg)
		if msg != nil {
			m.focus = enums.FocusSend
			m.amount.Focus()
		}
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "r":
		return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenLoading} }
	case "p", "P":
		m.focus = enums.FocusContacts
		// m.statusMessage = "Recipient picker active"
		return m, nil
	case "s":
		m.focus = enums.FocusSimulate
	case "S":
		return m.beginSend()
	case "N":
		m.networkDialog = networkDialog.New()
		m.networkDialog.Visible = true
		return m, nil
	case "n":
		m.networkPopup = networksPopup.New()
		m.networkPopup.Visible = true
		return m, nil
	}

	switch msg.String() {
	case "e", "E":
		m.focus = enums.FocusAmount
		m.amount.Focus()
		m.transaction.Reset()
		return m, nil
	case "t", "T":
		m.focus = enums.FocusTokens
		return m, nil
	case "x", "X":
		m.focus = enums.FocusTransactions
		return m, nil
	case "h", "left":
		m.focus = enums.FocusTokens
		return m, nil
	case "l", "right":
		m.focus = enums.FocusTransactions
		return m, nil

	}

	return m, nil
}
