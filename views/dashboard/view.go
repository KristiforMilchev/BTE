package dashboard

import (
	"strings"

	"bos/components"
	recipientPanel "bos/components/recipient_panel"
	"bos/constants"
	"bos/views"

	"github.com/charmbracelet/lipgloss"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

func (m *Model) View() string {
	base := views.RenderApp(
		m.width,
		m.height,
		m.focus,
		m.statusMessage,
		constants.RpcURL,
		func(width, height int) string {
			return m.renderDashboard(width, height)
		},
	)

	if m.networkDialog.Visible {
		return overlay.Composite(
			m.networkDialog.View(),
			base,
			overlay.Center,
			overlay.Center,
			0,
			0,
		)
	}

	if m.networkPopup.Visible {
		return overlay.Composite(
			m.networkPopup.View(),
			base,
			overlay.Left,
			overlay.Bottom,
			-2,
			-2,
		)
	}

	return base
}

func (m *Model) renderDashboard(width int, height int) string {
	outerGap := 2
	centerGap := 2

	availableWidth := width - outerGap - outerGap - centerGap
	if availableWidth < 40 {
		availableWidth = 40
	}

	leftWidth := availableWidth / 2
	rightWidth := availableWidth - leftWidth

	transferPanel := m.renderTransferPanelContent(leftWidth, height)
	assetsPanel := m.renderTokensPanel(rightWidth, height)

	return strings.Repeat(" ", outerGap) +
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			transferPanel,
			strings.Repeat(" ", centerGap),
			assetsPanel,
		) +
		strings.Repeat(" ", outerGap)
}

func (m *Model) renderTransferPanelContent(width int, height int) string {
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

	return PanelContentSized(width, height, body)
}

func (m *Model) renderTokensPanel(width int, height int) string {
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

	return PanelContentSized(width, height, body)
}

func PanelContentSized(width int, height int, body string) string {
	width = components.Max(8, width)

	style := components.PanelStyle.Width(width)

	if height > 0 {
		style = style.Height(components.Max(1, height-components.PanelStyle.GetVerticalFrameSize()))
	}

	return style.Render(body)
}

func PanelBodyWidth(width int) int {
	return components.Max(1, width-components.PanelStyle.GetHorizontalPadding())
}
