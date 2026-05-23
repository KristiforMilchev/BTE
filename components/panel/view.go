package panel

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

// Render creates a titled panel with the given content.
func Render(title string, width int, height int, content string) string {
	frameWidth := components.PanelStyle.GetHorizontalFrameSize()
	frameHeight := components.PanelStyle.GetVerticalFrameSize()

	innerWidth := components.Max(1, width-frameWidth)
	innerHeight := components.Max(1, height-frameHeight)

	body := strings.Join([]string{
		components.SectionTitle.
			Width(innerWidth).
			AlignHorizontal(lipgloss.Center).
			Render(title),
		"",
		content,
	}, "\n")

	return components.PanelSized(width, innerHeight, body)
}
