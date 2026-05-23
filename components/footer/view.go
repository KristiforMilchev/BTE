package footer

import (
	"bos/components"
	"bos/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderFooter(
	width int,
	statusMessage string,
	rpcURL string,
	address string,
) string {
	width = utils.SafeWidth(width)

	status := "Ledger Connected • RPC Online"

	if rpcURL != "" {
		status += " • " + rpcURL
	}

	if strings.TrimSpace(statusMessage) != "" {
		status += " • " + statusMessage
	}

	left := components.MutedText.Render(
		components.Truncate(status, width),
	)

	right := address

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)

	spacerWidth := max(0, width-leftWidth-rightWidth)

	spacer := strings.Repeat(" ", spacerWidth)

	row := left + spacer + right

	return components.Separator(width) + "\n" + row
}
