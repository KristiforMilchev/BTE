package views

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func RenderWalletPanel(s State, width int, height int) string {
	contentWidth := components.Max(28, width-components.PanelStyle.GetHorizontalFrameSize()-4)

	body := []string{
		components.SectionTitle.Render("Wallet"),
		"",
		components.KeyValue("Balance", walletBalance(s), contentWidth),
		components.KeyValue("Address", components.ShortAddress(s.Address), contentWidth),
		components.Separator(contentWidth),
		components.SectionTitle.Render("Assets"),
		"",
		renderTokenList(s, contentWidth),
		"",
	}

	return components.PanelSized(width, height, strings.Join(body, "\n"))
}

func renderTokenList(s State, width int) string {
	if len(s.Tokens) == 0 {
		return components.MutedText.Render("No assets loaded")
	}

	rows := make([]string, 0, len(s.Tokens))

	for i, token := range s.Tokens {
		active := s.Focus == FocusTokens &&
			s.SelectedToken == i

		selected := s.SelectedToken == i

		marker := "  "
		if active {
			marker = components.FocusMarker(true)
		}

		symbolStyle := components.Value

		if selected {
			symbolStyle = symbolStyle.
				Copy().
				Foreground(components.Accent)
		}

		amountSymbol := components.Truncate(
			token.Balance+" "+token.Symbol,
			width-4,
		)

		address := components.Truncate(
			token.Address,
			width-12,
		)

		badge := lipgloss.NewStyle().
			Foreground(components.Success).
			Bold(true).
			Render("✓")

		separator := components.MutedText.Render(
			strings.Repeat(
				"─",
				max(8, width-4),
			),
		)

		padding := strings.Repeat(
			" ",
			max(
				1,
				width-
					lipgloss.Width(address)-
					lipgloss.Width("✓")-
					8,
			),
		)

		line1 := marker +
			symbolStyle.Render(amountSymbol)

		line2 := "  " + separator

		line3 := "  " +
			components.MutedText.Render(address) +
			padding +
			badge

		block := strings.Join([]string{
			line1,
			line2,
			line3,
		}, "\n")

		rows = append(rows, block)
	}

	return strings.Join(rows, "\n\n")
}

func walletBalance(s State) string {
	if len(s.Tokens) == 0 {
		return "0 ETH"
	}
	return s.Tokens[0].Balance + " " + s.Tokens[0].Symbol
}
