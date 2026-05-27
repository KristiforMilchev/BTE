package dtos

import (
	"bos/types"
	"bos/utils"
	"strings"
)

func (r SimulatorPerformResponse) ToSimulationReport(session types.SimulationSession) types.SimulationReport {
	report := types.SimulationReport{
		Title:          simulationTitle(session),
		Status:         simulationStatus(r.Execution.Status),
		RiskLevel:      r.riskLevel(),
		GasEstimate:    session.Transaction.Gas,
		Summary:        r.simulationSummary(session),
		BalanceChanges: r.balanceChanges(session),
		TokenApprovals: tokenApprovals(r.ApprovalFindings),
		BytecodeChecks: []types.BytecodeCheck{
			{
				Address:    utils.FirstNonEmpty(r.Contract.Address, session.Address),
				IsContract: r.Contract.HasCode,
				Note:       codeHashNote(r.Contract.CodeHashSHA256),
			},
		},
		Calls: []types.ContractCall{
			{
				Depth:    0,
				From:     utils.FirstNonEmpty(session.Transaction.From, session.Caller),
				To:       utils.FirstNonEmpty(session.Transaction.To, session.Address),
				Function: r.simulationAction(session),
				Value:    simulationValue(session),
			},
		},
		FunctionCalls: functionCalls(r.FunctionCalls),
		Events:        r.events(),
		Warnings:      r.Warnings,
	}

	if r.Execution.TransactionHash != "" {
		report.Events = append(report.Events, types.EventLog{
			Contract: utils.FirstNonEmpty(r.Contract.Address, session.Address),
			Name:     "Transaction Hash",
			Details:  r.Execution.TransactionHash,
		})
	}
	return report
}

func simulationTitle(session types.SimulationSession) string {
	if session.AddressContract {
		return "Contract Simulation"
	}
	return "Transaction Simulation"
}

func (r SimulatorPerformResponse) simulationSummary(session types.SimulationSession) string {
	if r.Execution.Details != "" {
		return r.Execution.Details
	}
	if session.AddressContract {
		return "Function simulations completed before broadcast."
	}
	return "Native transfer simulation completed before broadcast."
}

func (r SimulatorPerformResponse) simulationAction(session types.SimulationSession) string {
	if session.AddressContract {
		return utils.FirstNonEmpty(r.Execution.Mode, "contract simulation")
	}
	return "regular wallet transfer"
}

func simulationValue(session types.SimulationSession) string {
	if session.Amount != "" {
		return session.TransferValue()
	}
	return utils.FirstNonEmpty(session.Transaction.Value, "0x0")
}

func functionCalls(calls []SimulatorFunctionCallReport) []types.SimulationFunctionCall {
	out := make([]types.SimulationFunctionCall, 0, len(calls))
	for _, call := range calls {
		out = append(out, types.SimulationFunctionCall{
			Function:     call.Function,
			Signature:    call.Signature,
			PassedData:   call.PassedData,
			Status:       simulationStatus(call.Status),
			Consequences: cleanConsequences(call.Consequences),
		})
	}
	return out
}

func cleanConsequences(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" ||
			strings.Contains(normalized, "debug_tracetransaction") ||
			strings.Contains(normalized, ".sha256=") {
			continue
		}
		out = append(out, value)
	}
	return out
}

func (r SimulatorPerformResponse) balanceChanges(session types.SimulationSession) []types.BalanceChange {
	changes := make([]types.BalanceChange, 0, 2)
	if r.Balances.CallerBefore != "" || r.Balances.CallerAfter != "" {
		changes = append(changes, types.BalanceChange{
			Asset:  session.Network.SymbolValue(),
			Before: r.Balances.CallerBefore,
			After:  r.Balances.CallerAfter,
			Delta:  balanceDelta(r.Balances.CallerBefore, r.Balances.CallerAfter),
		})
	}
	if r.Balances.AddressBefore != "" || r.Balances.AddressAfter != "" {
		changes = append(changes, types.BalanceChange{
			Asset:  recipientBalanceAsset(session),
			Before: r.Balances.AddressBefore,
			After:  r.Balances.AddressAfter,
			Delta:  balanceDelta(r.Balances.AddressBefore, r.Balances.AddressAfter),
		})
	}
	return changes
}

func recipientBalanceAsset(session types.SimulationSession) string {
	if session.AddressContract {
		return "Contract native balance"
	}
	return "Recipient native balance"
}

func tokenApprovals(findings []SimulatorApprovalFinding) []types.TokenApproval {
	approvals := make([]types.TokenApproval, 0, len(findings))
	for _, finding := range findings {
		approvals = append(approvals, types.TokenApproval{
			Token:   utils.FirstNonEmpty(finding.Type, "approval"),
			Spender: utils.FirstNonEmpty(finding.Selector, "unknown"),
			Amount:  finding.Description,
			Risk:    utils.FirstNonEmpty(finding.Severity, "medium"),
		})
	}
	return approvals
}

func (r SimulatorPerformResponse) events() []types.EventLog {
	events := make([]types.EventLog, 0, len(r.Contract.StateChanges)+1)
	for _, change := range r.Contract.StateChanges {
		events = append(events, types.EventLog{
			Contract: r.Contract.Address,
			Name:     "State Change",
			Details:  change,
		})
	}
	return events
}

func (r SimulatorPerformResponse) riskLevel() string {
	level := "Low"
	for _, finding := range r.ApprovalFindings {
		switch strings.ToLower(finding.Severity) {
		case "critical":
			return "Critical"
		case "high":
			level = "High"
		case "medium":
			if level == "Low" {
				level = "Medium"
			}
		}
	}
	if r.Execution.Status != "" && r.Execution.Status != "0x1" && !strings.EqualFold(r.Execution.Status, "success") {
		if level == "Low" {
			level = "Medium"
		}
	}
	return level
}

func simulationStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "0x1", "1", "success", "passed":
		return "passed"
	case "0x0", "0", "failed", "reverted":
		return "reverted"
	case "":
		return "simulated"
	default:
		return status
	}
}

func codeHashNote(codeHash string) string {
	if codeHash == "" {
		return "Contract bytecode proof captured."
	}
	return "Bytecode proof sha256=" + codeHash
}

func balanceDelta(before string, after string) string {
	if before == "" || after == "" {
		return ""
	}
	if before == after {
		return "0x0"
	}
	return before + " -> " + after
}
