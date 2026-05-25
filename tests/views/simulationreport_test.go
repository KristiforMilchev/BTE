package views_test

import (
	"strings"
	"testing"

	"bos/types"
	"bos/views/simulationreport"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSimulationReportViewDoesNotPanicWithSparseReport(t *testing.T) {
	model := simulationreport.New(types.SimulationReportPayload{
		Report: types.SimulationReport{
			Title:     "Transaction Simulation",
			Status:    "Waiting for signature",
			RiskLevel: "Pending",
			BalanceChanges: []types.BalanceChange{
				{Asset: "ETH", Before: "1.00", After: "0.99", Delta: "-0.01"},
			},
			Calls: []types.ContractCall{
				{From: "0x1111111111111111111111111111111111111111", To: "0x2222222222222222222222222222222222222222", Function: "regular wallet transfer", Value: "0.01 ETH"},
			},
		},
	})

	next, _ := model.Update(tea.WindowSizeMsg{Width: 140, Height: 40})
	rendered := next.View()

	if !strings.Contains(rendered, "Transaction Simulation") {
		t.Fatalf("rendered view does not contain title: %q", rendered)
	}
	if strings.Contains(rendered, "Checked Bytecode") {
		t.Fatalf("wallet report should not render checked bytecode section: %q", rendered)
	}
}
