package loading

import (
	"bos/interfaces"
	"bos/repositories"
	"bos/views"
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/common"
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

	var address *common.Address
	ledgerAddress := m.account()
	lastRepoAddress := m.getAccountFromRepository()

	if ledgerAddress == nil && lastRepoAddress == nil {
		return views.WalletLoadedMsg{Err: errors.New("Account not configured")}
	}

	if ledgerAddress == nil {
		address = lastRepoAddress
	}

	networkBalance, err := m.network.Balance(*address)
	if err != nil {
		return views.WalletLoadedMsg{Err: err}
	}

	return views.WalletLoadedMsg{
		Address: address.Hex(),
		Balance: networkBalance.Balance,
		ChainID: networkBalance.ChainID.String(),
	}
}

func (m *Model) account() *common.Address {
	wallet, err := m.wallet.Account()
	if err != nil {
		log.Printf("Ledger not connected -> %s", err)
		return nil
	}
	
	return &wallet.Address
}

func (m *Model) getAccountFromRepository() *common.Address {
	exists, _ := m.accountsRepository.Account()
	return exists
}
