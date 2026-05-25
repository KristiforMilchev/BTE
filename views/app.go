package views

import (
	activeWalletLabel "bos/components/active_wallet_label"
	"bos/components/footer"
	"bos/components/header"
	"bos/enums"
	"bos/utils"
	"strings"

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
	wallet := activeWalletLabel.New()
	headerView := header.RenderHeader(width, focus, strings.TrimSpace(rpcURL) != "")
	footerView := footer.RenderFooter(width, statusMessage, rpcURL, wallet.View())

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
