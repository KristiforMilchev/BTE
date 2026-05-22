package views

import (
	"bos/components"
	"bos/components/footer"
	"bos/components/header"
	"bos/enums"
	"bos/utils"

	"github.com/charmbracelet/lipgloss"
)

func RenderApp(width int, height int, focus enums.FocusArea, statusMessage string, rpcURL string, content string) string {
	width = utils.SafeWidth(width)
	if height <= 0 {
		height = 30
	}

	return components.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			header.RenderHeader(width, focus),
			content,
			footer.RenderFooter(width, statusMessage, rpcURL),
		),
	)
}
