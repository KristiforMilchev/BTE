package dashboard

import (
	"bos/components/amount"
	"bos/components/contacts"
	tokenlist "bos/components/token_list"
	transactionPreview "bos/components/transaction_preview"
	"bos/enums"
	"bos/interfaces"
	"bos/types"
	"bos/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Wallet        interfaces.IWallet
	Address       string
	Balance       string
	ChainID       string
	statusMessage string
}

type Model struct {
	wallet interfaces.IWallet

	width  int
	height int

	address string
	balance string
	chainID string

	statusMessage string

	focus enums.FocusArea

	contacts    *contacts.Model
	transaction *transactionPreview.Model
	amount      *amount.Model
	tokenList   *tokenlist.Model
}

func New(config Config) *Model {

	return &Model{
		wallet:      config.Wallet,
		address:     config.Address,
		balance:     config.Balance,
		chainID:     config.ChainID,
		focus:       enums.FocusSend,
		contacts:    contacts.NewContacts(),
		amount:      amount.New(),
		tokenList:   tokenlist.New(),
		transaction: transactionPreview.New(6),
	}
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) beginSend() (tea.Model, tea.Cmd) {
	amount := strings.TrimSpace(m.amount.Value())
	if amount == "" {
		m.statusMessage = "Enter an amount before sending"
		m.focus = enums.FocusAmount
		return m, nil
	}

	if _, err := utils.ParseEtherToWei(amount); err != nil {
		m.statusMessage = "Invalid amount: " + err.Error()
		m.focus = enums.FocusAmount
		return m, nil
	}
	if !m.tokenList.SelectedAsset().Native {
		m.statusMessage = "Token transfer signing is not integrated yet"
		return m, nil
	}
	if !common.IsHexAddress(m.contacts.SelectedRecipient().Address) {
		m.statusMessage = "Selected contact has an invalid address"
		m.focus = enums.FocusContacts
		return m, nil
	}

	draft := types.TxDraft{
		FromAddress: m.address, RecipientName: m.contacts.SelectedRecipient().Name, RecipientAddress: m.contacts.SelectedRecipient().Address,
		Amount: amount, Asset: m.tokenList.SelectedAsset(), EstimatedFee: m.transaction.EstimatedFee(),
		SimulationStatus: m.transaction.SimulationStatus(), RiskLevel: m.transaction.RiskLevel(),
	}
	return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenConfirm, Payload: draft} }
}
