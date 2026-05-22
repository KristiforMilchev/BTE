package sending

import (
	"strings"

	"bos/components"
	"bos/views"

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

	body := lipgloss.NewStyle().Padding(2, 4).Render(components.Panel(56, content))
	return views.RenderApp(m.width, m.height, views.FocusSend, "Waiting for Ledger approval", "", body)
}
