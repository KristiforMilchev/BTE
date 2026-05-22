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
	"bos/views/sending"
	"bos/views/sent"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	current tea.Model

	wallet   interfaces.IWallet
	network  interfaces.INetwork
	register repositories.RepositoryRegister

	width  int
	height int

	dashboard *dashboard.Model
}

func New(wallet interfaces.IWallet, network interfaces.INetwork, register repositories.RepositoryRegister) *Model {
	return &Model{
		wallet:   wallet,
		network:  network,
		current:  loading.New(wallet, network, register.Accounts),
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
			m.current = errorview.New(types.ErrorPayload{
				Title:   "Transaction Failed",
				Message: utils.ErrorMessage(msg.Err),
				Return:  enums.ScreenDashboard,
			})
			return m, resizeCmd(m.width, m.height)
		}

		m.current = sent.New(types.SentPayload{TxHash: msg.TxHash})
		return m, resizeCmd(m.width, m.height)
	}

	next, cmd := m.current.Update(msg)
	m.current = next
	return m, cmd
}

func (m *Model) View() string {
	if m.current == nil {
		return ""
	}
	return m.current.View()
}

func (m *Model) navigate(msg types.NavigateMsg) tea.Cmd {
	switch msg.Screen {
	case enums.ScreenLoading:
		m.current = loading.New(m.wallet, m.network, m.register.Accounts)
		return tea.Batch(resizeCmd(m.width, m.height), m.current.Init())

	case enums.ScreenDashboard:
		if m.dashboard != nil {
			m.current = m.dashboard
			return resizeCmd(m.width, m.height)
		}
		m.current = loading.New(m.wallet, m.network, m.register.Accounts)
		return tea.Batch(resizeCmd(m.width, m.height), m.current.Init())

	case enums.ScreenConfirm:
		draft, ok := msg.Payload.(types.TxDraft)
		if !ok {
			m.current = errorview.New(types.ErrorPayload{Title: "Invalid View Payload", Message: "missing transaction draft", Return: enums.ScreenDashboard})
			return resizeCmd(m.width, m.height)
		}
		m.current = confirm.New(m.wallet, draft)
		return resizeCmd(m.width, m.height)

	case enums.ScreenSending:
		m.current = sending.New()
		return resizeCmd(m.width, m.height)

	case enums.ScreenSent:
		payload, _ := msg.Payload.(types.SentPayload)
		m.current = sent.New(payload)
		return resizeCmd(m.width, m.height)

	case enums.ScreenError:
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
