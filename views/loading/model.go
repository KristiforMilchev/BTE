package loading

import (
	"bos/interfaces"
	"bos/repositories"
	"bos/types"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	wallet             interfaces.IWallet
	network            interfaces.INetwork
	accountsRepository repositories.AccountsRepository
	width              int
	height             int
}

func New(wallet interfaces.IWallet, network interfaces.INetwork, accounts repositories.AccountsRepository) *Model {
	return &Model{wallet: wallet, network: network, accountsRepository: accounts}
}

func (m *Model) Init() tea.Cmd {
	return m.loadWallet
}

func (m *Model) loadWallet() tea.Msg {

	wallet, err := m.wallet.Account()
	if err != nil {
		log.Printf("Ledger not connected -> %s", err)
		return types.WalletLoadedMsg{Err: err}
	}

	networkBalance, err := m.network.Balance(*wallet)
	if err != nil {
		return types.WalletLoadedMsg{Err: err}
	}

	return types.WalletLoadedMsg{
		Address: wallet.Hex(),
		Balance: networkBalance.Balance,
		ChainID: networkBalance.ChainID.String(),
	}
}
