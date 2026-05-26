package confirm

import (
	"fmt"

	"bos/enums"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "y", "enter", " ":
			return m, tea.Batch(
				func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenSending} },
				m.sendTransaction(),
			)
		case "n", "esc":
			return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenDashboard} }
		}
	}
	return m, nil
}

func (m *Model) sendTransaction() tea.Cmd {
	recipient := m.draft.RecipientAddress
	amount := m.draft.Amount

	return func() tea.Msg {
		txHash, err := m.wallet.SendTransaction(recipient, &amount, nil)
		if err != nil {
			return types.SendFinishedMsg{Err: err}
		}
		if txHash == nil {
			return types.SendFinishedMsg{Err: fmt.Errorf("wallet returned an empty transaction hash")}
		}
		return types.SendFinishedMsg{TxHash: *txHash, Draft: m.draft}
	}
}
