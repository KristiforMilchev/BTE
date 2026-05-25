package header

import (
	"bos/components"
	networkStatus "bos/components/network_status"
	"bos/enums"
	"bos/utils"

	"github.com/charmbracelet/lipgloss"
)

func RenderHeader(width int, focus enums.FocusArea, networkOnline bool) string {
	width = utils.SafeWidth(width)
	sideGap := 2
	minBadgeGap := 1
	contentWidth := components.Max(1, width-sideGap-sideGap)

	left := networkStatus.Render(networkOnline)
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
