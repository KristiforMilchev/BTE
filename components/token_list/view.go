package tokenlist

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Render(width int) string {
	if len(m.tokens) == 0 {
		return components.MutedText.Render("No assets loaded")
	}

	innerWidth := components.Max(8, width)

	rows := make([]string, 0, len(m.tokens))

	for i, token := range m.tokens {
		selected := m.selectedToken == i

		symbolStyle := components.Value
		if selected {
			symbolStyle = symbolStyle.Copy().Foreground(components.Accent)
		}

		amountSymbol := components.Truncate(
			token.Balance+" "+token.Symbol,
			innerWidth,
		)

		address := components.Truncate(
			token.Address,
			innerWidth-4,
		)

		addr := components.MutedText.Render(address)

		badge := lipgloss.NewStyle().
			Foreground(components.Success).
			Bold(true).
			Render("✓")

		gap := components.Max(
			0,
			innerWidth-2-lipgloss.Width(addr)-lipgloss.Width(badge),
		)

		bottomLine := "  " + addr + strings.Repeat(" ", gap) + badge

		block := strings.Join([]string{
			symbolStyle.Width(innerWidth).MaxWidth(innerWidth).Render(amountSymbol),
			components.Separator(innerWidth),
			bottomLine,
		}, "\n")

		rows = append(rows, block)
	}

	return strings.Join(rows, "\n\n")
}

func (m *Model) View() string {
	return m.Render(80)
}

func (m *Model) ViewWidth(width int) string {
	return m.Render(width)
}
