package dtos

type SimulatorBalanceSnapshot struct {
	CallerBefore  string `json:"callerBefore"`
	CallerAfter   string `json:"callerAfter,omitempty"`
	AddressBefore string `json:"addressBefore"`
	AddressAfter  string `json:"addressAfter,omitempty"`
}
