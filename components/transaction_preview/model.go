package transactionPreview

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	fee              string
	risk             string
	estimatedFee     string
	riskLevel        string
	width            int
	simulationStatus string
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) EstimatedFee() string {
	return m.estimatedFee
}

func (m *Model) RiskLevel() string {
	return m.riskLevel
}

func (m *Model) SimulationStatus() string {
	return m.simulationStatus
}

func (m *Model) Reset() {
	m.simulationStatus = "Not Run"
	m.riskLevel = "—"
	m.estimatedFee = "—"
}



func New(width int) *Model {
	return &Model{
		fee:              "-",
		risk:             "-",
		width:            width,
		simulationStatus: "Not Run",
		riskLevel:        "—",
		estimatedFee:     "—",
	}
}
