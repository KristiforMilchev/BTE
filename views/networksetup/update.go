package networksetup

import (
	networkDialog "bos/components/network_dialog"
	"bos/enums"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return m.dialog.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case networkDialog.SubmittedMsg:
		network, err := m.register.Network.NetworkByRpc(&msg.Network.RPC)
		if err == nil && network != nil {
			m.network.Change(network)
		}
		return m, func() tea.Msg {
			return types.NavigateMsg{Screen: enums.ScreenLoading}
		}

	case networkDialog.CancelledMsg:
		return m, nil
	}

	next, cmd := m.dialog.Update(msg)
	m.dialog = next
	return m, cmd
}
