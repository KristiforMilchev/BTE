package views

type Screen int

const (
	ScreenLoading Screen = iota
	ScreenDashboard
	ScreenConfirm
	ScreenSending
	ScreenSent
	ScreenError
)

type FocusArea int

const (
	FocusAmount FocusArea = iota
	FocusContacts
	FocusTokens
	FocusSimulate
	FocusSend
)

type Contact struct {
	Name    string
	Address string
}

type Token struct {
	Symbol   string
	Name     string
	Balance  string
	Address  string
	Native   bool
	Verified bool
}

type State struct {
	Screen Screen
	Focus  FocusArea

	Width  int
	Height int

	Address string
	Balance string
	ChainID string
	Err     string
	TxHash  string

	AmountInput string
	AmountValue string

	SelectedToken   int
	SelectedContact int

	Tokens   []Token
	Contacts []Contact

	SimulationStatus string
	RiskLevel        string
	EstimatedFee     string
	StatusMessage    string
	RpcURL           string
}
