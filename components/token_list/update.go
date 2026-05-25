package tokenlist

import (
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {
	case "left", "h":

	case "right", "l":

	case "up", "k":
		m.moveSelection(-1)
		return nil, nil
	case "down", "j":
		m.moveSelection(1)
		return nil, nil
	case "enter", "space":
		// return m.activateFocusedItem()
	}

	return m, nil
}

func (m *Model) moveSelection(delta int) {
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
}

func (m Model) SelectedAsset() types.Token {
	if len(m.tokens) == 0 {
		return types.Token{Symbol: "ETH", Balance: "0", Native: true}
	}
	if m.selectedToken < 0 || m.selectedToken >= len(m.tokens) {
		return m.tokens[0]
	}
	return m.tokens[m.selectedToken]
}
