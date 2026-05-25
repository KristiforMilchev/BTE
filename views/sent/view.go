package sent

import (
	"strings"

	"bos/components"
	"bos/enums"
	"bos/views"

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

	body := lipgloss.NewStyle().Padding(2, 4).Render(components.Panel(66, content))
	return views.RenderApp(m.width, m.height, enums.FocusSend, "Transaction broadcast", func(width, height int) string {
		return body
	})
}
