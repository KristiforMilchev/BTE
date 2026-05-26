package amount

import (
	"bos/components"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

func (m *Model) View() string {
	amount := strings.TrimSpace(m.amountInput.Value())
	if amount == "" {
		amount = "0.00"
	}

	amountStyle := components.HeroAmount
	if m.active {
		amountStyle = amountStyle.Foreground(components.Accent)
	}

	fig := figure.NewFigure(amount, "big", true)

	amountLine := amountStyle.Render(strings.TrimRight(fig.String(), "\n"))
	symbolLine := components.SectionTitle.Render(m.token.Symbol)

	width := max(lipgloss.Width(amountLine), lipgloss.Width(symbolLine))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(width, lipgloss.Center, amountLine),
		lipgloss.PlaceHorizontal(width, lipgloss.Center, symbolLine),
	)
}
