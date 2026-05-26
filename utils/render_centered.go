package utils

import (
	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func RenderCentered(width int, height int, title string, body string) string {
	width = SafeWidth(width)
	if height <= 0 {
		height = 30
	}

	content := components.SectionTitle.Render(title) + "\n\n" + body
	box := components.Panel(64, content)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}
