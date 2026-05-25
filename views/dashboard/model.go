package dashboard

import (
	"bos/components/amount"
	"bos/components/contacts"
	networkDialog "bos/components/network_dialog"
	networksPopup "bos/components/network_popup"
	tokenlist "bos/components/token_list"
	transactionPreview "bos/components/transaction_preview"
	transactionsComponent "bos/components/transactions"
	"bos/di"
	"bos/enums"
	"bos/interfaces"
	"bos/types"
	"bos/utils"
	"log"
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

	statusMessage string

	focus enums.FocusArea

	contacts     *contacts.Model
	transaction  *transactionPreview.Model
	amount       *amount.Model
	tokenList    *tokenlist.Model
	transactions *transactionsComponent.Model

	networkDialog *networkDialog.Model
	networkPopup  *networksPopup.Model
}

func New(config Config) *Model {

	return &Model{
		wallet:        config.Wallet,
		focus:         enums.FocusSend,
		contacts:      contacts.NewContacts(),
		amount:        amount.New(),
		tokenList:     tokenlist.New(),
		transactions:  transactionsComponent.New(sampleTransactions()),
		transaction:   transactionPreview.New(6),
		networkDialog: networkDialog.New(),
		networkPopup:  networksPopup.New(),
	}
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) onNetworkChanged() {
	m.tokenList = tokenlist.New()

}

func sampleTransactions() []types.Transaction {
	return []types.Transaction{
		{
			To:     "0x1111111111111111111111111111111111111111",
			Block:  "19420431",
			TxHash: "0x8f4f0af2042a7a7ff36b9a2a90131b1981a3d6907f3fb62d3b93e16ee13f6e2a",
			Amount: "0.00126092 ETH",
		},
		{
			To:     "0x2222222222222222222222222222222222222222",
			Block:  "19420412",
			TxHash: "0x73ceac921ea905abbfa29d45314d4dcb5f9d0d858df580d20d7608a55dbf8f21",
			Amount: "0.045 ETH",
		},
		{
			To:     "0x3333333333333333333333333333333333333333",
			Block:  "19420390",
			TxHash: "0x4da4619adf0093b43e37cf370d6f873d48d936f3a41c6fb1bb58883f78195832",
			Amount: "1.250 ETH",
		},
		{
			To:     "0x4444444444444444444444444444444444444444",
			Block:  "19420344",
			TxHash: "0xd2a78a19737a4b79b8e0156f8b9e384b202ddda78ce0a07d74185970aca88c1d",
			Amount: "0.004 ETH",
		},
		{
			To:     "0x5555555555555555555555555555555555555555",
			Block:  "19420281",
			TxHash: "0xadd6e94a8214c3f5e3730d7d10e0cb977b7a7a72e60c4f43039a862f314923c7",
			Amount: "0.083 ETH",
		},
	}
}

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

	account, err := di.GetWallet().Account()
	if err != nil {
		log.Printf("Can't start transaction account is nil -> %s", err)
		return m, nil
	}

	draft := types.TxDraft{
		FromAddress: account.Hex(), RecipientName: m.contacts.SelectedRecipient().Name, RecipientAddress: m.contacts.SelectedRecipient().Address,
		Amount: amount, Asset: m.tokenList.SelectedAsset(), EstimatedFee: m.transaction.EstimatedFee(),
		SimulationStatus: m.transaction.SimulationStatus(), RiskLevel: m.transaction.RiskLevel(),
	}
	return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenConfirm, Payload: draft} }
}
