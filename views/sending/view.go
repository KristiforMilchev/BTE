package sending

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	content := strings.Join([]string{
		components.SectionTitle.Render("Waiting for Ledger"),
		"",
		"Approve the transaction on your device.",
		"",
		components.MutedText.Render("Do not disconnect the Ledger."),
	}, "\n")

	return lipgloss.NewStyle().Padding(1, 2).Render(components.Panel(56, content))
}
