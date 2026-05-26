package contractinteraction

import (
	"bos/enums"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenDashboard} }
		case "k", "up":
			m.moveSelection(-1)
			return m, nil
		case "j", "down":
			m.moveSelection(1)
			return m, nil
		case "enter", "s":
			report := fakeContractSimulation(m.address, m.functions[m.selected])
			return m, func() tea.Msg {
				return types.NavigateMsg{
					Screen: enums.ScreenSimulationReport,
					Payload: types.SimulationReportPayload{
						Report: report,
						Return: "contract",
					},
				}
			}
		}
	}

	return m, nil
}

func (m *Model) moveSelection(delta int) {
	if len(m.functions) == 0 {
		return
	}
	m.selected += delta
	if m.selected < 0 {
		m.selected = len(m.functions) - 1
	}
	if m.selected >= len(m.functions) {
		m.selected = 0
	}
}

func fakeContractSimulation(address string, function types.ContractFunction) types.SimulationReport {
	return types.SimulationReport{
		Title:       "Contract Simulation",
		Status:      "Simulated",
		RiskLevel:   "Medium",
		GasEstimate: "84,210",
		Summary:     "Fake fork execution for " + function.Signature + ". External API integration will replace this payload.",
		BalanceChanges: []types.BalanceChange{
			{Asset: "ETH", Before: "1.0000", After: "0.9975", Delta: "-0.0025"},
		},
		TokenApprovals: []types.TokenApproval{
			{Token: "MOCK", Spender: "0x2222222222222222222222222222222222222222", Amount: "unlimited", Risk: "high"},
		},
		Calls: []types.ContractCall{
			{Depth: 0, From: "0xWallet", To: address, Function: function.Signature, Value: "0"},
			{Depth: 1, From: address, To: "0xTokenVault", Function: "updateAccounting(address,uint256)", Value: "0"},
		},
		Events: []types.EventLog{
			{Contract: address, Name: "CallSimulated", Details: function.Name + " completed on fork"},
		},
		Warnings: []string{
			"Function mutability is " + function.Mutability + ".",
			"ABI/function list is fake until external discovery is wired.",
		},
	}
}
