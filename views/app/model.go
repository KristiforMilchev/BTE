package app

import (
	"bos/interfaces"
	"bos/repositories"
	"bos/views"
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

	case views.WalletLoadedMsg:
		if msg.Err != nil {
			m.current = errorview.New(views.ErrorPayload{
				Title:   "Blockcert",
				Message: views.ErrorMessage(msg.Err),
				Return:  views.ScreenLoading,
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

	case views.NavigateMsg:
		return m, m.navigate(msg)

	case views.SendFinishedMsg:
		if msg.Err != nil {
			m.current = errorview.New(views.ErrorPayload{
				Title:   "Transaction Failed",
				Message: views.ErrorMessage(msg.Err),
				Return:  views.ScreenDashboard,
			})
			return m, resizeCmd(m.width, m.height)
		}

		m.current = sent.New(views.SentPayload{TxHash: msg.TxHash})
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

func (m *Model) navigate(msg views.NavigateMsg) tea.Cmd {
	switch msg.Screen {
	case views.ScreenLoading:
		m.current = loading.New(m.wallet, m.network, m.register.Accounts)
		return tea.Batch(resizeCmd(m.width, m.height), m.current.Init())

	case views.ScreenDashboard:
		if m.dashboard != nil {
			m.current = m.dashboard
			return resizeCmd(m.width, m.height)
		}
		m.current = loading.New(m.wallet, m.network, m.register.Accounts)
		return tea.Batch(resizeCmd(m.width, m.height), m.current.Init())

	case views.ScreenConfirm:
		draft, ok := msg.Payload.(views.TxDraft)
		if !ok {
			m.current = errorview.New(views.ErrorPayload{Title: "Invalid View Payload", Message: "missing transaction draft", Return: views.ScreenDashboard})
			return resizeCmd(m.width, m.height)
		}
		m.current = confirm.New(m.wallet, draft)
		return resizeCmd(m.width, m.height)

	case views.ScreenSending:
		m.current = sending.New()
		return resizeCmd(m.width, m.height)

	case views.ScreenSent:
		payload, _ := msg.Payload.(views.SentPayload)
		m.current = sent.New(payload)
		return resizeCmd(m.width, m.height)

	case views.ScreenError:
		payload, ok := msg.Payload.(views.ErrorPayload)
		if !ok {
			payload = views.ErrorPayload{Title: "Blockcert", Message: "unknown error", Return: views.ScreenDashboard}
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
