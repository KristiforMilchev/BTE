package amount

import (
	"bos/components"
	"bos/layout"
	"strings"
)

func (m Model) View() string {

	amountStyle := components.HeroAmount
	if m.active {
		amountStyle = amountStyle.Foreground(components.Accent)
	}

	amount := strings.TrimSpace(m.amountInput.Value())
	amountDisplay := amount

	body :=
		layout.Column(
			
			layout.Expanded(1, layout.WidgetFunc(func(ctx layout.Context) string {
				return amountStyle.Render(amountDisplay)
			})),
			layout.Expanded(1, layout.WidgetFunc(func(ctx layout.Context) string {
				return components.SectionTitle.Render(m.token.Symbol)
			})),
		)

	return body.Render(layout.Context{

		Constraints: layout.Constraints{
			Height: 50,
		},

	})
}
