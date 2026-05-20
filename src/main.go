package main

import (
	"fmt"
	"bos/constants"
	"bos/di"
	"bos/interfaces"
	"bos/utils"
	"math/big"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ethereum/go-ethereum/common"
)

type screen int

const (
	screenLoading screen = iota
	screenHome
	screenRecipient
	screenAmount
	screenConfirm
	screenSending
	screenSent
	screenError
)

type ledgerResultMsg struct {
	Address string
	Balance string
	ChainID *big.Int
	Err     error
}

type sendResultMsg struct {
	TxHash string
	Err    error
}

type Model struct {
	screen screen

	address string
	balance string
	chainID string
	err     string
	txHash  string

	recipientInput textinput.Model
	amountInput    textinput.Model

	recipient string
	amount    string

	wallet  interfaces.IWallet
	network interfaces.INetwork
}

func main() {
	di.SetupDependencies()

	recipient := textinput.New()
	recipient.Placeholder = "0x..."
	recipient.CharLimit = 42
	recipient.Width = 46

	amount := textinput.New()
	amount.Placeholder = "0.01"
	amount.CharLimit = 32
	amount.Width = 32

	p := tea.NewProgram(&Model{
		screen:         screenLoading,
		recipientInput: recipient,
		amountInput:    amount,
		network:        di.GetNetwork(),
		wallet:         di.GetWallet(),
	})

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (m *Model) Init() tea.Cmd {
	wallet, err := m.wallet.Account()
	if err != nil {
		m.screen = screenError
		m.err = err.Error()
		return nil
	}

	networkBalance, err := m.network.Balance(wallet.Address)
	if err != nil {
		m.screen = screenError
		m.err = err.Error()
		return nil
	}

	m.address = wallet.Address.Hex()
	m.balance = networkBalance.Balance
	m.chainID = networkBalance.ChainID.String()
	m.screen = screenHome

	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			m.screen = screenLoading
			m.err = ""
			m.txHash = ""

			m.Init()
			return m, nil
		}

		switch m.screen {
		case screenHome:
			if msg.String() == "s" {
				m.recipientInput.SetValue("")
				m.amountInput.SetValue("")
				m.recipientInput.Focus()
				m.screen = screenRecipient
				return m, textinput.Blink
			}

		case screenRecipient:
			if msg.String() == "enter" {
				to := strings.TrimSpace(m.recipientInput.Value())
				if !common.IsHexAddress(to) {
					m.screen = screenError
					m.err = "invalid recipient address"
					return m, nil
				}

				m.recipient = common.HexToAddress(to).Hex()
				m.amountInput.Focus()
				m.screen = screenAmount
				return m, textinput.Blink
			}

			var cmd tea.Cmd
			m.recipientInput, cmd = m.recipientInput.Update(msg)
			return m, cmd

		case screenAmount:
			if msg.String() == "enter" {
				amount := strings.TrimSpace(m.amountInput.Value())
				if _, err := utils.ParseEtherToWei(amount); err != nil {
					m.screen = screenError
					m.err = err.Error()
					return m, nil
				}

				m.amount = amount
				m.screen = screenConfirm
				return m, nil
			}

			var cmd tea.Cmd
			m.amountInput, cmd = m.amountInput.Update(msg)
			return m, cmd

		case screenConfirm:
			if msg.String() == "y" {
				m.screen = screenSending
				txHash, err := m.wallet.SendTransaction(m.recipient, &m.amount, nil)
				if err != nil {
					m.screen = screenError
					m.err = err.Error()
					return m, nil
				}
				m.txHash = *txHash
				m.screen = screenSent
				return m, nil
			}

			if msg.String() == "n" || msg.String() == "esc" {
				m.screen = screenHome
				return m, nil
			}

		case screenSent:
			if msg.String() == "enter" {
				m.screen = screenLoading
				m.txHash = ""

				m.Init()
				return m, nil
			}

		case screenError:
			if msg.String() == "enter" {
				m.screen = screenLoading
				m.err = ""

				m.Init()
				return m, nil
			}
		}
	}

	return m, nil
}

func (m *Model) View() string {
	switch m.screen {
	case screenLoading:
		return `
Blockcert Ledger Test

Connecting to Ledger...

Requirements:
- Ledger plugged in
- Device unlocked
- Ethereum app open
- Ledger Live closed

Press q to quit.
`

	case screenHome:
		return fmt.Sprintf(`
Blockcert Ledger Test

RPC:
%s

Chain ID:
%s

Ledger Address:
%s

Balance:
%s ETH

Options:
s = send transaction
r = refresh
q = quit
`, constants.RpcURL, m.chainID, m.address, m.balance)

	case screenRecipient:
		return fmt.Sprintf(`
Send Transaction

From:
%s

Enter recipient address:

%s

Press enter to continue.
Press q to quit.
`, m.address, m.recipientInput.View())

	case screenAmount:
		return fmt.Sprintf(`
Send Transaction

From:
%s

To:
%s

Enter amount in ETH:

%s

Press enter to continue.
Press q to quit.
`, m.address, m.recipient, m.amountInput.View())

	case screenConfirm:
		return fmt.Sprintf(`
Confirm Transaction

From:
%s

To:
%s

Amount:
%s ETH

The Ledger device will ask you to confirm the transaction.

Press y to sign and broadcast.
Press n to cancel.
`, m.address, m.recipient, m.amount)

	case screenSending:
		return `
Sending Transaction

Check your Ledger device and approve the transaction.

Press q to quit.
`

	case screenSent:
		return fmt.Sprintf(`
Transaction Sent

Hash:
%s

Press enter to refresh balance.
Press q to quit.
`, m.txHash)

	case screenError:
		return fmt.Sprintf(`
Blockcert Ledger Test

Error:
%s

Press enter to retry.
Press q to quit.
`, m.err)
	}

	return ""
}
