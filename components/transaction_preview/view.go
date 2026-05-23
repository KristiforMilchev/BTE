package transactionPreview

import (
	"bos/components"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	return strings.Join([]string{
		components.SectionTitle.Render("Preview"),
		"",
		components.KeyValue("Fee", m.estimatedFee, m.width),
		components.KeyValue("Risk", riskLabel(m.riskLevel), m.width),
		components.KeyValue("Simulation", m.simulationStatus, m.width),
	}, "\n")

}

func riskLabel(risk string) string {
	switch strings.ToLower(risk) {
	case "low":
		return lipgloss.NewStyle().Foreground(components.Success).Bold(true).Render(risk)
	case "medium":
		return lipgloss.NewStyle().Foreground(components.Warning).Bold(true).Render(risk)
	case "high", "critical":
		return lipgloss.NewStyle().Foreground(components.Danger).Bold(true).Render(risk)
	default:
		return risk
	}
}
