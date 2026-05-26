package types

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
