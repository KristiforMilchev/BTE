package dashboard

import (
	"strings"

	"bos/components"
	recipientPanel "bos/components/recipient_panel"
	"bos/constants"
	"bos/views"

	"github.com/charmbracelet/lipgloss"
	"github.com/consensys/gnark-crypto/utils"
)

func (m *Model) View() string {
	return views.RenderApp(
		m.width,
		m.height,
		m.focus,
		m.statusMessage,
		constants.RpcURL,
		m.renderDashboard(),
	)
}

func (m *Model) renderDashboard() string {
	frame := components.PanelStyle.GetHorizontalFrameSize()

	bodyWidth := m.width - frame*2
	leftWidth := bodyWidth / 2
	rightWidth := bodyWidth - leftWidth

	transferPanel := m.renderTransferPanelContent(leftWidth)
	assetsPanel := m.renderTokensPanel(rightWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top, transferPanel, assetsPanel)
}

func (m *Model) renderTransferPanelContent(width int) string {
	bodyWidth := PanelBodyWidth(width)

	asset := m.tokenList.SelectedAsset()
	m.amount.SetSymbol(asset)

	recipient := m.contacts.SelectedRecipient()

	body := strings.Join([]string{
		lipgloss.PlaceHorizontal(bodyWidth, lipgloss.Center, m.amount.View()),
		"",
		recipientPanel.Render(recipient, bodyWidth),
		"",
		m.transaction.View(),
		"",
		components.Separator(bodyWidth),
		components.SectionTitle.Render("Contacts"),
		"",
		m.contacts.ViewWidth(bodyWidth),
	}, "\n")

	return PanelContentSized(width, m.height-4, body)
}

func (m *Model) renderTokensPanel(width int) string {
	bodyWidth := PanelBodyWidth(width)

	body := strings.Join([]string{
		components.SectionTitle.
			Width(bodyWidth).
			MaxWidth(bodyWidth).
			AlignHorizontal(lipgloss.Center).
			Render("Assets"),
		"",
		m.tokenList.ViewWidth(bodyWidth),
	}, "\n")

	return PanelContentSized(width, m.height-4, body)
}
func PanelContentSized(width int, height int, body string) string {
	width = utils.Max(8, width)

	style := components.PanelStyle.Width(width)

	if height > 0 {
		style = style.Height(utils.Max(1, height))
	}

	return style.Render(body)
}
func PanelBodyWidth(width int) int {
	return utils.Max(1, width-components.PanelStyle.GetHorizontalPadding())
}
