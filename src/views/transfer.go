package views

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func RenderTransferPanel(s State, width int, height int) string {
	asset := selectedAsset(s)
	recipient := selectedRecipient(s)

	innerWidth := components.Max(32, width-components.PanelStyle.GetHorizontalFrameSize()-4)

	amount := strings.TrimSpace(s.AmountValue)
	amountDisplay := amount
	if amountDisplay == "" {
		amountDisplay = "0.00"
	}

	body := strings.Join([]string{
		renderAmountHero(amountDisplay, asset.Symbol, s.Focus == FocusAmount, innerWidth),
		"",
		renderRecipientBlock(recipient, innerWidth),
		"",
		renderPreviewBlock(s, innerWidth),
		"",
		components.Separator(innerWidth),
		components.SectionTitle.Render("Contacts"),
		"",
		RenderContacts(s, innerWidth),
	}, "\n")

	return components.PanelSized(width, height, body)
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

func renderRecipientBlock(recipient Contact, width int) string {
	return strings.Join([]string{
		components.SectionTitle.Render("Recipient"),
		components.Value.Render(components.Truncate(recipient.Name, width)),
		components.MutedText.Render(components.ShortAddress(recipient.Address)),
	}, "\n")
}

func renderPreviewBlock(s State, width int) string {
	return strings.Join([]string{
		components.SectionTitle.Render("Preview"),
		"",
		components.KeyValue("Fee", s.EstimatedFee, width),
		components.KeyValue("Risk", riskLabel(s.RiskLevel), width),
		components.KeyValue("Simulation", s.SimulationStatus, width),
	}, "\n")
}
