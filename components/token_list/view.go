package tokenlist

import (
	"strings"

	"bos/components"
	"bos/types"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Render(width int) string {
	return m.RenderHeight(width, 0)
}

func (m *Model) RenderHeight(width int, height int) string {
	if len(m.tokens) == 0 {
		return components.MutedText.Render("No assets loaded")
	}

	innerWidth := components.Max(28, width)
	assetWidth, balanceWidth, typeWidth := tokenTableColumns(innerWidth)

	header := components.Value.Render(
		tokenCell("ASSET", assetWidth) + "  " +
			tokenCell("BALANCE", balanceWidth) + "  " +
			tokenCell("TYPE", typeWidth),
	)

	rows := []string{
		"  " + header,
		components.Separator(innerWidth),
	}

	visibleRows := len(m.tokens)
	if height > 0 {
		visibleRows = components.Max(1, height-len(rows))
	}
	m.ensureSelectedVisible(visibleRows)

	end := components.Min(len(m.tokens), m.offset+visibleRows)
	for i := m.offset; i < end; i++ {
		token := m.tokens[i]
		selected := m.selectedToken == i

		marker := "  "
		rowStyle := components.Value
		if selected {
			marker = components.FocusMarker(true)
			rowStyle = rowStyle.Copy().Foreground(components.Accent)
		}

		row := tokenCell(token.Symbol, assetWidth) + "  " +
			tokenCell(token.Balance, balanceWidth) + "  " +
			tokenCell(tokenType(token), typeWidth)

		rows = append(rows, marker+rowStyle.Render(row))
	}

	return strings.Join(rows, "\n")
}

func (m *Model) View() string {
	return m.Render(80)
}

func (m *Model) ViewWidth(width int) string {
	return m.Render(width)
}

func (m *Model) ViewWidthHeight(width int, height int) string {
	return m.RenderHeight(width, height)
}

func (m *Model) ensureSelectedVisible(visibleRows int) {
	if visibleRows < 1 {
		visibleRows = 1
	}
	if m.selectedToken < m.offset {
		m.offset = m.selectedToken
	}
	if m.selectedToken >= m.offset+visibleRows {
		m.offset = m.selectedToken - visibleRows + 1
	}
	if maxOffset := components.Max(0, len(m.tokens)-visibleRows); m.offset > maxOffset {
		m.offset = maxOffset
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func tokenTableColumns(width int) (int, int, int) {
	typeWidth := components.Clamp(width/5, 6, 10)
	balanceWidth := components.Clamp(width/3, 10, 18)
	assetWidth := width - 2 - balanceWidth - 2 - typeWidth
	if assetWidth < 8 {
		assetWidth = 8
	}

	return assetWidth, balanceWidth, typeWidth
}

func tokenCell(value string, width int) string {
	value = components.Truncate(value, width)
	padding := width - lipgloss.Width(value)
	if padding < 0 {
		padding = 0
	}

	return value + strings.Repeat(" ", padding)
}

func tokenType(token types.Token) string {
	if token.Native {
		return "native"
	}
	if token.Verified {
		return "verified"
	}
	return "token"
}
