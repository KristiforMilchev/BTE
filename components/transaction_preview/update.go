package transactionPreview

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) Update() (tea.Msg, tea.Cmd) {
	return m, nil
}

func (m *Model) SetSimulationStatus(data string) {
	m.simulationStatus = data
}

func (m *Model) SetRiskLevel(data string) {
	m.riskLevel = data
	m.risk = data
}

func (m *Model) SetEstimatedFee(data string) {
	m.estimatedFee = data
}
