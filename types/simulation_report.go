package types

import "strings"

func (r SimulationReport) IsEmpty() bool {
	return r.Title == "" &&
		r.Status == "" &&
		r.RiskLevel == "" &&
		r.GasEstimate == "" &&
		r.Summary == "" &&
		len(r.BalanceChanges) == 0 &&
		len(r.TokenApprovals) == 0 &&
		len(r.BytecodeChecks) == 0 &&
		len(r.Calls) == 0 &&
		len(r.Events) == 0 &&
		len(r.Warnings) == 0
}

func (r SimulationReport) WithFallbacks(session SimulationSession) SimulationReport {
	if r.Title == "" {
		r.Title = "Transaction Simulation"
	}
	if r.Status == "" {
		r.Status = "Simulated"
	}
	if r.RiskLevel == "" {
		r.RiskLevel = "Low"
	}
	if r.Summary == "" {
		r.Summary = "Signed transaction was tested before broadcast."
	}
	if len(r.BalanceChanges) == 0 {
		r.BalanceChanges = []BalanceChange{
			{Asset: session.AssetSymbol(), Before: session.BalanceBefore, After: session.BalanceAfter, Delta: "-" + session.Amount},
		}
	}
	if len(r.BytecodeChecks) == 0 {
		r.BytecodeChecks = []BytecodeCheck{session.BytecodeCheck()}
	}
	if len(r.Calls) == 0 {
		r.Calls = session.PerformedActions()
	}
	r.Warnings = meaningfulWarnings(r.Warnings)
	return r
}

func (s SimulationSession) AssetSymbol() string {
	if s.Asset != "" {
		return s.Asset
	}
	return s.Network.SymbolValue()
}

func (s SimulationSession) TransferValue() string {
	return s.Amount + " " + s.AssetSymbol()
}

func (s SimulationSession) PerformedActions() []ContractCall {
	if !s.AddressContract {
		return []ContractCall{
			{
				Depth:    0,
				From:     s.Caller,
				To:       s.Address,
				Function: "regular wallet transfer",
				Value:    s.TransferValue(),
			},
		}
	}

	return []ContractCall{
		{Depth: 0, From: s.Caller, To: s.Address, Function: "contract call", Value: s.TransferValue()},
	}
}

func (s SimulationSession) BytecodeCheck() BytecodeCheck {
	if !s.AddressContract {
		return BytecodeCheck{
			Address:       s.Address,
			IsContract:    false,
			RuntimeHex:    "0x",
			RuntimeBinary: "0b0",
			Note:          "No runtime bytecode returned for this address.",
		}
	}

	runtimeHex := firstNonEmpty(s.RuntimeHex, "0x")
	runtimeBinary := firstNonEmpty(s.RuntimeBinary, "0b1")
	return BytecodeCheck{
		Address:       s.Address,
		IsContract:    true,
		RuntimeHex:    runtimeHex,
		RuntimeBinary: runtimeBinary,
		Note:          "Contract bytecode proof captured.",
	}
}

func (n Network) SymbolValue() string {
	if n.Symbol == nil || *n.Symbol == "" {
		return "ETH"
	}
	return *n.Symbol
}

func meaningfulWarnings(warnings []string) []string {
	meaningful := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		normalized := strings.ToLower(strings.TrimSpace(warning))
		if normalized == "" {
			continue
		}
		if normalized == "no authentication or rate limit checks are applied" {
			continue
		}
		if normalized == "signed transaction was executed only on the ganache fork and was not broadcast to the upstream network" {
			continue
		}
		meaningful = append(meaningful, warning)
	}
	return meaningful
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
