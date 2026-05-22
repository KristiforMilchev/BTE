package views

import "fmt"

type NavigateMsg struct {
	Screen  Screen
	Payload any
}

type WalletLoadedMsg struct {
	Address string
	Balance string
	ChainID string
	Err     error
}

type SendFinishedMsg struct {
	TxHash string
	Err    error
}

type ErrorPayload struct {
	Title   string
	Message string
	Return  Screen
}

type SentPayload struct {
	TxHash string
}

type TxDraft struct {
	FromAddress      string
	RecipientName    string
	RecipientAddress string
	Amount           string
	Asset            Token
	EstimatedFee     string
	SimulationStatus string
	RiskLevel        string
}

func ErrorMessage(err error) string {
	if err == nil {
		return "unknown error"
	}
	return fmt.Sprintf("%s", err)
}
