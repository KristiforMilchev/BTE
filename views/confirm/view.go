package confirm

import (
	"strings"

	"bos/components"
	"bos/views"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	contentWidth := 58
	amount := strings.TrimSpace(m.draft.Amount)

	content := strings.Join([]string{
		components.SectionTitle.Render("Confirm Transaction"),
		"",
		components.KeyValue("From", components.ShortAddress(m.draft.FromAddress), contentWidth),
		components.KeyValue("To", m.draft.RecipientName, contentWidth),
		components.KeyValue("Address", components.ShortAddress(m.draft.RecipientAddress), contentWidth),
		components.KeyValue("Amount", amount+" "+m.draft.Asset.Symbol, contentWidth),
		components.KeyValue("Fee", m.draft.EstimatedFee, contentWidth),
		components.KeyValue("Simulation", m.draft.SimulationStatus, contentWidth),
		components.KeyValue("Risk Level", riskLabel(m.draft.RiskLevel), contentWidth),
		"",
		components.MutedText.Render("Only approve on your Ledger if these details are correct."),
		"",
		components.Button("y  Sign on Ledger", true) + " " + components.Button("n  Cancel", false),
	}, "\n")

	body := lipgloss.NewStyle().Padding(2, 4).Render(components.Panel(66, content))
	return views.RenderApp(m.width, m.height, views.FocusSend, "Confirm transaction", "", body)
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
