package sent

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	content := strings.Join([]string{
		components.SuccessText.Render("Transaction Sent"),
		"",
		components.Label.Render("Hash"),
		components.Value.Render(components.Truncate(m.payload.TxHash, 58)),
		"",
		components.MutedText.Render("Press enter to return to dashboard."),
	}, "\n")

	return lipgloss.NewStyle().Padding(1, 2).Render(components.Panel(66, content))
}
