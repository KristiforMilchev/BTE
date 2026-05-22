package dashboard

import (
	"bos/enums"
	"bos/interfaces"
	"bos/types"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Config struct {
	Wallet  interfaces.IWallet
	Address string
	Balance string
	ChainID string
}

type Model struct {
	wallet interfaces.IWallet

	width  int
	height int

	address string
	balance string
	chainID string

	amountInput textinput.Model

	focus enums.FocusArea

	selectedToken   int
	selectedContact int

	tokens   []types.Token
	contacts []types.Contact

	simulationStatus string
	riskLevel        string
	estimatedFee     string
	statusMessage    string
}

func New(config Config) *Model {
	amount := textinput.New()
	amount.Placeholder = "0.01"
	amount.CharLimit = 32
	amount.Width = 0
	amount.Focus()

	tokens := fakeTokens()
	if len(tokens) > 0 {
		tokens[0].Balance = config.Balance
	}

	return &Model{
		wallet:           config.Wallet,
		address:          config.Address,
		balance:          config.Balance,
		chainID:          config.ChainID,
		amountInput:      amount,
		focus:            enums.FocusAmount,
		selectedToken:    0,
		selectedContact:  0,
		tokens:           tokens,
		contacts:         fakeContacts(),
		simulationStatus: "Not Run",
		riskLevel:        "—",
		estimatedFee:     "—",
		statusMessage:    "Wallet loaded",
	}
}

func (m *Model) Init() tea.Cmd { return nil }
