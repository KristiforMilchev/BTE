package networksetup

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		components.SectionTitle.Render("Setup First Network"),
		components.MutedText.Render("Add a network before opening the dashboard."),
		"",
		m.dialog.View(),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		strings.TrimRight(content, "\n"),
	)
}
