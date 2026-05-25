package header

import (
	"bos/components"
	networkStatus "bos/components/network_status"
	"bos/di"
	"bos/enums"
	"bos/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderHeader(width int, focus enums.FocusArea) string {
	width = utils.SafeWidth(width)
	sideGap := 2
	minBadgeGap := 1
	contentWidth := components.Max(1, width-sideGap-sideGap)

	left := networkStatus.Render(hasNetworkRPC())
	right := components.MutedText.Render(activeHelp(focus))

	headerHeight := lipgloss.Height(left)
	right = lipgloss.NewStyle().
		Width(lipgloss.Width(right)).
		Height(headerHeight).
		AlignVertical(lipgloss.Center).
		Render(right)

	spacerWidth := max(minBadgeGap, contentWidth-lipgloss.Width(left)-lipgloss.Width(right))
	spacer := lipgloss.NewStyle().
		Width(spacerWidth).
		Height(headerHeight).
		Render("")
	rightPadding := lipgloss.NewStyle().
		Width(sideGap).
		Height(headerHeight).
		Render("")
	leftPadding := lipgloss.NewStyle().
		Width(sideGap).
		Height(headerHeight).
		Render("")

	line := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPadding,
		left,
		spacer,
		right,
		rightPadding,
	)

	return line + "\n" + components.Separator(width)
}

func hasNetworkRPC() bool {
	networkProvider := di.GetNetwork()
	if networkProvider == nil {
		return false
	}
	currentNetwork := networkProvider.Network()
	if currentNetwork.Rpc == nil {
		return false
	}
	return strings.TrimSpace(*currentNetwork.Rpc) != ""
}
