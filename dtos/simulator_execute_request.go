package dtos

import "bos/types"

type SimulatorExecuteRequest struct {
	SimulationID       string                              `json:"simulationId"`
	Session            types.SimulationSession             `json:"session"`
	SignedTx           string                              `json:"signedTx,omitempty"`
	SignedTransaction  string                              `json:"signedTransaction,omitempty"`
	RawTransaction     string                              `json:"rawTransaction,omitempty"`
	SignedTransactions []types.SignedSimulationTransaction `json:"signedTransactions,omitempty"`
}
