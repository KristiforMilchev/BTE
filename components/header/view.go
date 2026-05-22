package header

import (
	"bos/components"
	"bos/enums"
	"bos/utils"

	"github.com/charmbracelet/lipgloss"
)

func RenderHeader(width int, focus enums.FocusArea) string {
	width = utils.SafeWidth(width)

	left := components.Value.Render("Blockcert")
	right := components.MutedText.Render(activeHelp(focus))

	spacerWidth := components.Max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-2)
	line := left + lipgloss.NewStyle().Width(spacerWidth).Render("") + right

	return line + "\n" + components.Separator(width)
}
