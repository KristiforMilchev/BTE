package dtos

type SimulatorExecutionReport struct {
	Mode              string `json:"mode"`
	Broadcasted       bool   `json:"broadcasted"`
	TransactionHash   string `json:"transactionHash,omitempty"`
	Status            string `json:"status"`
	Details           string `json:"details"`
	ForkBackendNeeded bool   `json:"forkBackendNeeded"`
}
