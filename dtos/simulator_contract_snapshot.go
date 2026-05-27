package dtos

type SimulatorContractSnapshot struct {
	Address        string   `json:"address"`
	HasCode        bool     `json:"hasCode"`
	CodeHashSHA256 string   `json:"codeHashSha256,omitempty"`
	StateChanges   []string `json:"stateChanges,omitempty"`
}
