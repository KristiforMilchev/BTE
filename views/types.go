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
