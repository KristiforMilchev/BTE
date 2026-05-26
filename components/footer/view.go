package footer

import (
	"bos/components"
	"bos/di"
	"bos/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderFooter(
	width int,
	statusMessage string,
	address string,
) string {
	width = utils.SafeWidth(width)
	sideGap := 2
	contentWidth := components.Max(1, width-sideGap-sideGap)

	status := "Ledger Connected • "

	if networkStatus := activeNetworkStatus(); networkStatus != "" {
		status += networkStatus
	}

	if strings.TrimSpace(statusMessage) != "" {
		status += " • " + statusMessage
	}

	right := address
	rightWidth := lipgloss.Width(right)

	left := components.MutedText.Render(
		components.Truncate(status, components.Max(1, contentWidth-rightWidth)),
	)

	leftWidth := lipgloss.Width(left)
	spacerWidth := max(0, contentWidth-leftWidth-rightWidth)

	spacer := strings.Repeat(" ", spacerWidth)

	row := strings.Repeat(" ", sideGap) + left + spacer + right + strings.Repeat(" ", sideGap)

	return components.Separator(width) + "\n" + row
}

func activeNetworkStatus() string {
	networkProvider := di.GetNetwork()
	if networkProvider == nil {
		return ""
	}
	network := networkProvider.Network()
	if network.Rpc == nil || strings.TrimSpace(*network.Rpc) == "" {
		return ""
	}

	name := ""
	if network.Name != nil {
		name = *network.Name
	}
	return "Network • " + name + " • RPC Online • " + *network.Rpc
}
