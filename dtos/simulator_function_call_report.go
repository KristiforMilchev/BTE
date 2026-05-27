package dtos

type SimulatorFunctionCallReport struct {
	Function     string   `json:"function"`
	Signature    string   `json:"signature"`
	PassedData   string   `json:"passedData"`
	Status       string   `json:"status"`
	Consequences []string `json:"consequences"`
}
