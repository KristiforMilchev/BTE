package dtos

import (
	"bos/types"
	"time"
)

type SimulatorBeginRequest struct {
	Network   string                   `json:"network"`
	Address   string                   `json:"address"`
	Caller    string                   `json:"caller"`
	Functions []types.ContractFunction `json:"functions,omitempty"`
}

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

type SimulatorExecuteRequest struct {
	SimulationID       string                              `json:"simulationId"`
	Session            types.SimulationSession             `json:"session"`
	SignedTx           string                              `json:"signedTx,omitempty"`
	SignedTransaction  string                              `json:"signedTransaction,omitempty"`
	RawTransaction     string                              `json:"rawTransaction,omitempty"`
	SignedTransactions []types.SignedSimulationTransaction `json:"signedTransactions,omitempty"`
}

type SimulatorExecuteResponse struct {
	Report types.SimulationReport `json:"report"`
}

type SimulatorPerformResponse struct {
	SimulationID         string                        `json:"simulationId"`
	Network              string                        `json:"network"`
	RawTransactionSHA256 string                        `json:"rawTransactionSha256"`
	Balances             SimulatorBalanceSnapshot      `json:"balances"`
	ApprovalFindings     []SimulatorApprovalFinding    `json:"approvalFindings"`
	Contract             SimulatorContractSnapshot     `json:"contract"`
	Execution            SimulatorExecutionReport      `json:"execution"`
	FunctionCalls        []SimulatorFunctionCallReport `json:"functionCalls,omitempty"`
	Warnings             []string                      `json:"warnings,omitempty"`
	PerformedAt          time.Time                     `json:"performedAt"`
}

type SimulatorBalanceSnapshot struct {
	CallerBefore  string `json:"callerBefore"`
	CallerAfter   string `json:"callerAfter,omitempty"`
	AddressBefore string `json:"addressBefore"`
	AddressAfter  string `json:"addressAfter,omitempty"`
}

type SimulatorApprovalFinding struct {
	Type        string `json:"type"`
	Selector    string `json:"selector"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type SimulatorContractSnapshot struct {
	Address        string   `json:"address"`
	HasCode        bool     `json:"hasCode"`
	CodeHashSHA256 string   `json:"codeHashSha256,omitempty"`
	StateChanges   []string `json:"stateChanges,omitempty"`
}

type SimulatorExecutionReport struct {
	Mode              string `json:"mode"`
	Broadcasted       bool   `json:"broadcasted"`
	TransactionHash   string `json:"transactionHash,omitempty"`
	Status            string `json:"status"`
	Details           string `json:"details"`
	ForkBackendNeeded bool   `json:"forkBackendNeeded"`
}

type SimulatorFunctionCallReport struct {
	Function     string   `json:"function"`
	Signature    string   `json:"signature"`
	PassedData   string   `json:"passedData"`
	Status       string   `json:"status"`
	Consequences []string `json:"consequences"`
}
