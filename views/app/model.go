package app

import (
	"bos/enums"
	"bos/interfaces"
	"bos/repositories"
	"bos/types"
	"bos/utils"
	"bos/views/confirm"
	"bos/views/dashboard"
	"bos/views/errorview"
	"bos/views/loading"
	"bos/views/networksetup"
	"bos/views/sending"
	"bos/views/sent"

	tea "github.com/charmbracelet/bubbletea"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

type Model struct {
	current tea.Model
	modal   tea.Model

	wallet   interfaces.IWallet
	network  interfaces.INetwork
	register repositories.RepositoryRegister

	width  int
	height int

	dashboard *dashboard.Model
}

func New(wallet interfaces.IWallet, network interfaces.INetwork, register repositories.RepositoryRegister) *Model {
	current := tea.Model(loading.New(wallet, network, register.Accounts))
	if !hasActiveNetwork(network.Network()) {
		current = networksetup.New(network, register)
	}

	return &Model{
		wallet:   wallet,
		network:  network,
		current:  current,
		register: register,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.current.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.current != nil {
			next, cmd := m.current.Update(msg)
			m.current = next
			if m.modal != nil {
				modal, modalCmd := m.modal.Update(msg)
				m.modal = modal
				return m, tea.Batch(cmd, modalCmd)
			}
			return m, cmd
		}

	case types.WalletLoadedMsg:
		if msg.Err != nil {
			m.current = errorview.New(types.ErrorPayload{
				Title:   "Blockcert",
				Message: utils.ErrorMessage(msg.Err),
				Return:  enums.ScreenLoading,
			})
			return m, nil
		}

		m.dashboard = dashboard.New(dashboard.Config{
			Wallet:  m.wallet,
			Address: msg.Address,
			Balance: msg.Balance,
			ChainID: msg.ChainID,
		})
		m.current = m.dashboard
		return m, resizeCmd(m.width, m.height)

	case types.NavigateMsg:
		return m, m.navigate(msg)

	case types.SendFinishedMsg:
		if msg.Err != nil {
			m.modal = nil
			m.current = errorview.New(types.ErrorPayload{
				Title:   "Transaction Failed",
				Message: utils.ErrorMessage(msg.Err),
				Return:  enums.ScreenDashboard,
			})
			return m, resizeCmd(m.width, m.height)
		}

		if m.dashboard != nil {
			m.dashboard.OnTransactionSent()
		}
		m.modal = sent.New(types.SentPayload{TxHash: msg.TxHash})
		return m, resizeCmd(m.width, m.height)
	}

	if m.modal != nil {
		next, cmd := m.modal.Update(msg)
		m.modal = next
		return m, cmd
	}

	next, cmd := m.current.Update(msg)
	m.current = next
	return m, cmd
}

func (m *Model) View() string {
	if m.current == nil {
		return ""
	}

	base := m.current.View()
	if m.modal == nil {
		return base
	}

	return overlay.Composite(
		m.modal.View(),
		base,
		overlay.Center,
		overlay.Center,
		0,
		0,
	)
}

func (m *Model) navigate(msg types.NavigateMsg) tea.Cmd {
	switch msg.Screen {
	case enums.ScreenLoading:
		m.modal = nil
		m.current = loading.New(m.wallet, m.network, m.register.Accounts)
		return tea.Batch(resizeCmd(m.width, m.height), m.current.Init())

	case enums.ScreenDashboard:
		m.modal = nil
		if m.dashboard != nil {
			m.current = m.dashboard
			return resizeCmd(m.width, m.height)
		}
		m.current = loading.New(m.wallet, m.network, m.register.Accounts)
		return tea.Batch(resizeCmd(m.width, m.height), m.current.Init())

	case enums.ScreenConfirm:
		draft, ok := msg.Payload.(types.TxDraft)
		if !ok {
			m.modal = nil
			m.current = errorview.New(types.ErrorPayload{Title: "Invalid View Payload", Message: "missing transaction draft", Return: enums.ScreenDashboard})
			return resizeCmd(m.width, m.height)
		}
		m.modal = confirm.New(m.wallet, draft)
		return resizeCmd(m.width, m.height)

	case enums.ScreenSending:
		m.modal = sending.New()
		return resizeCmd(m.width, m.height)

	case enums.ScreenSent:
		payload, _ := msg.Payload.(types.SentPayload)
		m.modal = sent.New(payload)
		return resizeCmd(m.width, m.height)

	case enums.ScreenError:
		m.modal = nil
		payload, ok := msg.Payload.(types.ErrorPayload)
		if !ok {
			payload = types.ErrorPayload{Title: "Blockcert", Message: "unknown error", Return: enums.ScreenDashboard}
		}
		m.current = errorview.New(payload)
		return resizeCmd(m.width, m.height)
	}

	return nil
}

func resizeCmd(width, height int) tea.Cmd {
	if width <= 0 || height <= 0 {
		return nil
	}
	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: width, Height: height}
	}
}

func hasActiveNetwork(network types.Network) bool {
	return network.Name != nil &&
		network.Rpc != nil &&
		network.Symbol != nil &&
		network.Chain != nil
}
