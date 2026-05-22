package tokenlist

import (
	"strings"

	"bos/components"
	"bos/layout"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Render(ctx layout.Context) string {
	width := ctx.Constraints.Width

	if len(m.tokens) == 0 {
		return components.MutedText.Render("No assets loaded")
	}

	rows := make([]string, 0, len(m.tokens))

	for i, token := range m.tokens {
		selected := m.selectedToken == i

		symbolStyle := components.Value
		if selected {
			symbolStyle = symbolStyle.Foreground(components.Accent)
		}

		amountSymbol := components.Truncate(token.Balance+" "+token.Symbol, width-4)
		address := components.Truncate(token.Address, width-12)

		badge := lipgloss.NewStyle().Foreground(components.Success).Bold(true).Render("✓")
		separator := components.MutedText.Render(strings.Repeat("─", max(8, width-4)))
		padding := strings.Repeat(" ", max(1, width-lipgloss.Width(address)-lipgloss.Width("✓")-8))

		block := strings.Join([]string{
			symbolStyle.Render(amountSymbol),
			"  " + separator,
			"  " + components.MutedText.Render(address) + padding + badge,
		}, "\n")

		rows = append(rows, block)
	}

	return strings.Join(rows, "\n\n")
}

func (m *Model) View() string {
	return m.Render(layout.Context{Constraints: layout.Constraints{Width: 80, Height: 24}})
}
