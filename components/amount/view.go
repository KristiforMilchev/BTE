package amount

import (
	"bos/components"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	amount := strings.TrimSpace(m.amountInput.Value())
	if amount == "" {
		amount = "0.00"
	}

	amountStyle := components.HeroAmount
	if m.active {
		amountStyle = amountStyle.Foreground(components.Accent)
	}

	amountLine := amountStyle.Render(amount)
	symbolLine := components.SectionTitle.Render(m.token.Symbol)

	width := max(lipgloss.Width(amountLine), lipgloss.Width(symbolLine))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(width, lipgloss.Center, amountLine),
		lipgloss.PlaceHorizontal(width, lipgloss.Center, symbolLine),
	)
}
