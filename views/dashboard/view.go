package dashboard

import (
	"strings"

	"bos/components"
	"bos/components/panel"
	"bos/constants"
	"bos/types"
	"bos/views"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	body := m.renderDashboard()
	return views.RenderApp(m.width, m.height, m.focus, m.statusMessage, constants.RpcURL, body)
}

func (m *Model) renderDashboard() string {
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	transferPanel := lipgloss.NewStyle().
		Width(leftWidth).
		Height(m.height - 4).
		Render(m.renderTransferPanelContent(leftWidth))

	assetsContent := panel.Render("Assets", rightWidth, m.height-4, m.tokenList.ViewWidth(rightWidth))

	assetsPanel := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height - 4).
		Render(assetsContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, transferPanel, assetsPanel)
}

func (m *Model) renderTransferPanelContent(width int) string {
	asset := m.tokenList.SelectedAsset()
	m.amount.SetSymbol(asset)
	recipient := m.contacts.SelectedRecipient()

	innerWidth := components.Max(32, width-components.PanelStyle.GetHorizontalFrameSize()-4)

	amount := strings.TrimSpace(m.amount.Value())
	amountDisplay := amount
	if amountDisplay == "" {
		amountDisplay = "0.00"
	}

	body := strings.Join([]string{
		lipgloss.PlaceHorizontal(innerWidth, lipgloss.Center, m.amount.View()),
		"",
		renderRecipientBlock(recipient, innerWidth),
		"",
		m.renderPreviewBlock(innerWidth),
		"",
		components.Separator(innerWidth),
		components.SectionTitle.Render("Contacts"),
		"",
		m.contacts.ViewWidth(innerWidth),
	}, "\n")

	return components.PanelSized(width, m.height-4, body)
}

func renderAmountHero(amount string, symbol string, active bool, width int) string {
	amountStyle := components.HeroAmount
	if active {
		amountStyle = amountStyle.Copy().Foreground(components.Accent)
	}

	amountLine := lipgloss.PlaceHorizontal(width, lipgloss.Center, amountStyle.Render(amount))
	symbolLine := lipgloss.PlaceHorizontal(width, lipgloss.Center, components.SectionTitle.Render(symbol))

	return "\n" + amountLine + "\n\n" + symbolLine + "\n"
}

func renderRecipientBlock(recipient types.Contact, width int) string {
	return strings.Join([]string{
		components.SectionTitle.Render("Recipient"),
		components.Value.Render(components.Truncate(recipient.Name, width)),
		components.MutedText.Render(components.ShortAddress(recipient.Address)),
	}, "\n")
}

func (m *Model) renderPreviewBlock(width int) string {
	return strings.Join([]string{
		components.SectionTitle.Render("Preview"),
		"",
		components.KeyValue("Fee", m.estimatedFee, width),
		components.KeyValue("Risk", riskLabel(m.riskLevel), width),
		components.KeyValue("Simulation", m.simulationStatus, width),
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
