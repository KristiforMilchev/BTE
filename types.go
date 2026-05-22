package main

import (
	"bos/interfaces"
	"bos/views"

	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	screen views.Screen
	focus  views.FocusArea

	width  int
	height int

	address string
	balance string
	chainID string
	err     string
	txHash  string

	amountInput textinput.Model

	selectedToken   int
	selectedContact int

	tokens   []views.Token
	contacts []views.Contact

	simulationStatus string
	riskLevel        string
	estimatedFee     string
	statusMessage    string

	wallet  interfaces.IWallet
	network interfaces.INetwork
}

type walletLoadedMsg struct {
	Address string
	Balance string
	ChainID string
	Err     error
}

type sendFinishedMsg struct {
	TxHash string
	Err    error
}
