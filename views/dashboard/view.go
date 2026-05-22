package dashboard

import (
	"strings"

	"bos/components"
	"bos/constants"
	"bos/enums"
	"bos/views"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	body := m.renderDashboard()
	return views.RenderApp(m.width, m.height, m.focus, m.statusMessage, constants.RpcURL, body)
}

func (m *Model) renderDashboard() string {
	outerWidth := safeWidth(m.width)
	gap := 3
	paddingX := 2
	paddingY := 1

	availableWidth := components.Max(90, outerWidth-(paddingX*2)-gap)
	leftWidth := availableWidth * 55 / 100
	rightWidth := availableWidth - leftWidth - gap

	availableHeight := components.Max(24, m.height-4-(paddingY*2))

	left := m.renderTransferPanel(leftWidth, availableHeight)
	right := m.renderWalletPanel(rightWidth, availableHeight)

	row := lipgloss.JoinHorizontal(lipgloss.Top, left, "   ", right)

	return lipgloss.NewStyle().
		Width(outerWidth).
		Padding(paddingY, paddingX).
		Render(row)
}

func (m *Model) renderTransferPanel(width int, height int) string {
	asset := m.selectedAsset()
	recipient := m.selectedRecipient()

	innerWidth := components.Max(32, width-components.PanelStyle.GetHorizontalFrameSize()-4)

	amount := strings.TrimSpace(m.amountInput.Value())
	amountDisplay := amount
	if amountDisplay == "" {
		amountDisplay = "0.00"
	}

	body := strings.Join([]string{
		renderAmountHero(amountDisplay, asset.Symbol, m.focus == enums.FocusAmount, innerWidth),
		"",
		renderRecipientBlock(recipient, innerWidth),
		"",
		m.renderPreviewBlock(innerWidth),
		"",
		components.Separator(innerWidth),
		components.SectionTitle.Render("Contacts"),
		"",
		m.renderContacts(innerWidth),
	}, "\n")

	return components.PanelSized(width, height, body)
}

func renderAmountHero(amount string, symbol string, active bool, width int) string {
	amountStyle := components.HeroAmount
	if active {
		amountStyle = amountStyle.Copy().Foreground(components.Accent)
	}

	amountLine := lipgloss.PlaceHorizontal(width, lipgloss.Center, amountStyle.Render(amount))
	symbolLine := lipgloss.PlaceHorizontal(width, lipgloss.Center, components.SectionTitle.Render(symbol))

	return "\n" + amountLine + "\n\n" + symbolLine + "\n"
}

func renderRecipientBlock(recipient views.Contact, width int) string {
	return strings.Join([]string{
		components.SectionTitle.Render("Recipient"),
		components.Value.Render(components.Truncate(recipient.Name, width)),
		components.MutedText.Render(components.ShortAddress(recipient.Address)),
	}, "\n")
}

func (m *Model) renderPreviewBlock(width int) string {
	return strings.Join([]string{
		components.SectionTitle.Render("Preview"),
		"",
		components.KeyValue("Fee", m.estimatedFee, width),
		components.KeyValue("Risk", riskLabel(m.riskLevel), width),
		components.KeyValue("Simulation", m.simulationStatus, width),
	}, "\n")
}

func (m *Model) renderContacts(width int) string {
	if len(m.contacts) == 0 {
		return components.MutedText.Render("No saved recipients")
	}

	rows := make([]string, 0, len(m.contacts))
	for i, contact := range m.contacts {
		active := m.focus == enums.FocusContacts && i == m.selectedContact
		selected := i == m.selectedContact

		nameStyle := components.Value
		marker := "  "
		if selected {
			nameStyle = nameStyle.Copy().Foreground(components.Accent)
		}
		if active {
			marker = components.FocusMarker(true)
		}

		row := marker + nameStyle.Render(components.Truncate(contact.Name, width-2)) + "\n" +
			"  " + components.MutedText.Render(components.ShortAddress(contact.Address))
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n\n")
}

func (m *Model) renderWalletPanel(width int, height int) string {
	contentWidth := components.Max(28, width-components.PanelStyle.GetHorizontalFrameSize()-4)

	body := []string{
		components.SectionTitle.Render("Wallet"),
		"",
		components.KeyValue("Balance", m.walletBalance(), contentWidth),
		components.KeyValue("Address", components.ShortAddress(m.address), contentWidth),
		components.Separator(contentWidth),
		components.SectionTitle.Render("Assets"),
		"",
		m.renderTokenList(contentWidth),
		"",
	}

	return components.PanelSized(width, height, strings.Join(body, "\n"))
}

func (m *Model) renderTokenList(width int) string {
	if len(m.tokens) == 0 {
		return components.MutedText.Render("No assets loaded")
	}

	rows := make([]string, 0, len(m.tokens))

	for i, token := range m.tokens {
		active := m.focus == enums.FocusTokens && m.selectedToken == i
		selected := m.selectedToken == i

		marker := "  "
		if active {
			marker = components.FocusMarker(true)
		}

		symbolStyle := components.Value
		if selected {
			symbolStyle = symbolStyle.Copy().Foreground(components.Accent)
		}

		amountSymbol := components.Truncate(token.Balance+" "+token.Symbol, width-4)
		address := components.Truncate(token.Address, width-12)

		badge := lipgloss.NewStyle().Foreground(components.Success).Bold(true).Render("✓")
		separator := components.MutedText.Render(strings.Repeat("─", max(8, width-4)))
		padding := strings.Repeat(" ", max(1, width-lipgloss.Width(address)-lipgloss.Width("✓")-8))

		block := strings.Join([]string{
			marker + symbolStyle.Render(amountSymbol),
			"  " + separator,
			"  " + components.MutedText.Render(address) + padding + badge,
		}, "\n")

		rows = append(rows, block)
	}

	return strings.Join(rows, "\n\n")
}

func (m *Model) walletBalance() string {
	if len(m.tokens) == 0 {
		return "0 ETH"
	}
	return m.tokens[0].Balance + " " + m.tokens[0].Symbol
}

func riskLabel(risk string) string {
	switch strings.ToLower(risk) {
	case "low":
		return lipgloss.NewStyle().Foreground(components.Success).Bold(true).Render(risk)
	case "medium":
		return lipgloss.NewStyle().Foreground(components.Warning).Bold(true).Render(risk)
	case "high", "critical":
		return lipgloss.NewStyle().Foreground(components.Danger).Bold(true).Render(risk)
	default:
		return risk
	}
}

func safeWidth(width int) int {
	if width < 100 {
		return 100
	}
	return width
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
