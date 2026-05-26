package transactions

import (
	"strings"

	"bos/components"
	"bos/types"

	"github.com/charmbracelet/lipgloss"
)

const transactionRowHeight = 1

func (m *Model) ViewWidthHeight(width int, height int) string {
	width = components.Max(24, width)
	height = components.Max(1, height)

	if len(m.transactions) == 0 {
		return components.MutedText.Render("No transactions")
	}

	visibleRows := components.Max(1, height-3)
	m.ensureSelectedVisible(visibleRows)

	end := components.Min(len(m.transactions), m.offset+visibleRows)
	rows := make([]string, 0, end-m.offset+2)
	rows = append(rows, m.header(width), components.Separator(width))
	for i := m.offset; i < end; i++ {
		rows = append(rows, m.row(m.transactions[i], width, i == m.selected))
	}

	return strings.Join(rows, "\n")
}

func (m *Model) View() string {
	return m.ViewWidthHeight(80, 20)
}

func (m *Model) header(width int) string {
	toWidth, amountWidth, hashWidth := transactionColumnWidths(width)
	return components.Label.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(toWidth).Render("To"),
			lipgloss.NewStyle().Width(amountWidth).Render("Amount"),
			lipgloss.NewStyle().Width(hashWidth).Render("Hash"),
		),
	)
}

func (m *Model) row(tx types.Transaction, width int, selected bool) string {
	toWidth, amountWidth, hashWidth := transactionColumnWidths(width)
	style := components.Value.Copy().Bold(false)
	if selected {
		style = style.Foreground(components.Accent).Bold(true)
	}

	return style.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(toWidth).Render(components.TruncateMiddle(tx.To, toWidth-1)),
			lipgloss.NewStyle().Width(amountWidth).Render(components.Truncate(tx.Amount, amountWidth-1)),
			lipgloss.NewStyle().Width(hashWidth).Render(components.TruncateMiddle(tx.TxHash, hashWidth)),
		),
	)
}

func transactionColumnWidths(width int) (int, int, int) {
	width = components.Max(32, width)
	amountWidth := components.Clamp(width/5, 8, 14)
	remaining := width - amountWidth
	toWidth := components.Max(10, remaining/2)
	hashWidth := components.Max(6, width-toWidth-amountWidth)
	return toWidth, amountWidth, hashWidth
}

func (m *Model) ensureSelectedVisible(visibleRows int) {
	if visibleRows < 1 {
		visibleRows = 1
	}
	if m.selected < m.offset {
		m.offset = m.selected
	}
	if m.selected >= m.offset+visibleRows {
		m.offset = m.selected - visibleRows + 1
	}
	if maxOffset := components.Max(0, len(m.transactions)-visibleRows); m.offset > maxOffset {
		m.offset = maxOffset
	}
	if m.offset < 0 {
		m.offset = 0
	}
}
