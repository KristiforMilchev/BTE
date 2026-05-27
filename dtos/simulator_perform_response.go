package dtos

import "time"

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
