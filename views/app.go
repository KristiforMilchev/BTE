package views

import (
	"bos/components/footer"
	"bos/components/header"
	"bos/enums"
	"bos/utils"

	"github.com/charmbracelet/lipgloss"
)

func RenderApp(
	width int,
	height int,
	focus enums.FocusArea,
	statusMessage string,
	rpcURL string,
	renderContent func(width int, height int) string,
) string {
	width = utils.SafeWidth(width)
	if height <= 0 {
		height = 30
	}

	headerView := header.RenderHeader(width, focus)
	footerView := footer.RenderFooter(width, statusMessage, rpcURL)

	bodyHeight := height -
		lipgloss.Height(headerView) -
		lipgloss.Height(footerView)

	if bodyHeight < 1 {
		bodyHeight = 1
	}

	body := renderContent(width, bodyHeight)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		body,
		footerView,
	)
}
