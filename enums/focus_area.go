package enums

type FocusArea int

const (
	FocusAmount FocusArea = iota
	FocusContacts
	FocusTokens
	FocusTransactions
	FocusSimulate
	FocusSend
)
