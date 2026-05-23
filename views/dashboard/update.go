package dashboard

import (
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
		msg, _ := m.tokenList.Update(msg)
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

	}

	return m, nil
}
