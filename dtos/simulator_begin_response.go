package dtos

import "bos/types"

type SimulatorBeginResponse struct {
	SimulationID    string                    `json:"simulationId"`
	ID              string                    `json:"id"`
	SessionID       string                    `json:"sessionId"`
	RPC             string                    `json:"rpc"`
	ClonedRPC       string                    `json:"clonedRpc"`
	ChainID         string                    `json:"chainId"`
	ChainIDAlt      string                    `json:"chainID"`
	Transaction     types.LedgerTransaction   `json:"transaction"`
	Transactions    []types.LedgerTransaction `json:"transactions"`
	Contract        SimulatorContractSnapshot `json:"contract"`
	AddressContract bool                      `json:"addressContract"`
	IsContract      bool                      `json:"isContract"`
	Bytecode        string                    `json:"bytecode"`
	RuntimeHex      string                    `json:"runtimeHex"`
	RuntimeBinary   string                    `json:"runtimeBinary"`
}
