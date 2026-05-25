package networkStatus

import (
	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func Render(online bool) string {
	label := " BLOCKCERT "
	color := components.Danger
	if online {
		label = " BLOCKCERT "
		color = components.Accent
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Bold(true).
		Render(label)
}
