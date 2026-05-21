package views

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func RenderConfirm(s State) string {
	asset := selectedAsset(s)
	recipient := selectedRecipient(s)
	amount := strings.TrimSpace(s.AmountValue)

	contentWidth := 58
	content := strings.Join([]string{
		components.SectionTitle.Render("Confirm Transaction"),
		"",
		components.KeyValue("From", components.ShortAddress(s.Address), contentWidth),
		components.KeyValue("To", recipient.Name, contentWidth),
		components.KeyValue("Address", components.ShortAddress(recipient.Address), contentWidth),
		components.KeyValue("Amount", amount+" "+asset.Symbol, contentWidth),
		components.KeyValue("Fee", s.EstimatedFee, contentWidth),
		components.KeyValue("Simulation", s.SimulationStatus, contentWidth),
		components.KeyValue("Risk Level", riskLabel(s.RiskLevel), contentWidth),
		"",
		components.MutedText.Render("Only approve on your Ledger if these details are correct."),
		"",
		components.Button("y  Sign on Ledger", true) + " " + components.Button("n  Cancel", false),
	}, "\n")

	return lipgloss.NewStyle().Padding(2, 4).Render(components.Panel(66, content))
}

func RenderSending(s State) string {
	content := strings.Join([]string{
		components.SectionTitle.Render("Waiting for Ledger"),
		"",
		"Approve the transaction on your device.",
		"",
		components.MutedText.Render("Do not disconnect the Ledger."),
	}, "\n")

	return lipgloss.NewStyle().Padding(2, 4).Render(components.Panel(56, content))
}

func RenderSent(s State) string {
	content := strings.Join([]string{
		components.SuccessText.Render("Transaction Sent"),
		"",
		components.Label.Render("Hash"),
		components.Value.Render(components.Truncate(s.TxHash, 58)),
		"",
		components.MutedText.Render("Press enter to return to dashboard."),
	}, "\n")

	return lipgloss.NewStyle().Padding(2, 4).Render(components.Panel(66, content))
}
