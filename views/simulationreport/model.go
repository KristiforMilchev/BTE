package simulationreport

import "bos/types"

type Model struct {
	payload types.SimulationReportPayload
	running bool
	width   int
	height  int
}

func New(payload types.SimulationReportPayload) *Model {
	if payload.Report.Title == "" && payload.Draft != nil {
		payload.Report = pendingReport(*payload.Draft)
	}
	return &Model{payload: payload}
}

type completedMsg struct {
	report types.SimulationReport
	err    error
}

type startMsg struct{}

func pendingReport(draft types.TxDraft) types.SimulationReport {
	return types.SimulationReport{
		Title:     "Transaction Simulation",
		Status:    "Waiting for signature",
		RiskLevel: "Pending",
		Summary:   "Review the Ledger prompt to sign a simulator-only transaction. The signed transaction will be sent to the simulator API, not broadcast.",
		BalanceChanges: []types.BalanceChange{
			{Asset: draft.Asset.Symbol, Before: draft.Asset.Balance, After: pendingBalanceAfter(draft), Delta: "-" + draft.Amount},
		},
		BytecodeChecks: []types.BytecodeCheck{
			{
				Address:       draft.RecipientAddress,
				IsContract:    false,
				RuntimeHex:    "0x",
				RuntimeBinary: "0b0",
				Note:          "No runtime bytecode detected in the pending fake check.",
			},
		},
		Calls: []types.ContractCall{
			{Depth: 0, From: draft.FromAddress, To: draft.RecipientAddress, Function: "regular wallet transfer", Value: draft.Amount + " " + draft.Asset.Symbol},
		},
	}
}
