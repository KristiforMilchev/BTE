package dtos

import "bos/types"

type SimulatorBeginRequest struct {
	Network   string                   `json:"network"`
	Address   string                   `json:"address"`
	Caller    string                   `json:"caller"`
	Functions []types.ContractFunction `json:"functions,omitempty"`
}
