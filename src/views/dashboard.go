package views

import (
	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func RenderDashboard(s State) string {
	outerWidth := s.Width
	gap := 3
	paddingX := 2
	paddingY := 1

	availableWidth := components.Max(90, outerWidth-(paddingX*2)-gap)
	leftWidth := availableWidth * 55 / 100
	rightWidth := availableWidth - leftWidth - gap

	availableHeight := components.Max(24, s.Height-4-(paddingY*2))

	left := RenderTransferPanel(s, leftWidth, availableHeight)
	right := RenderWalletPanel(s, rightWidth, availableHeight)

	row := lipgloss.JoinHorizontal(lipgloss.Top, left, "   ", right)

	return lipgloss.NewStyle().
		Width(outerWidth).
		Padding(paddingY, paddingX).
		Render(row)
}
