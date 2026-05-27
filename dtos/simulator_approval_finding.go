package dtos

type SimulatorApprovalFinding struct {
	Type        string `json:"type"`
	Selector    string `json:"selector"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}
