package transactions

import (
	"strings"

	"bos/components"
)

const transactionRowHeight = 5

func (m *Model) ViewWidthHeight(width int, height int) string {
	width = components.Max(12, width)
	height = components.Max(1, height)

	if len(m.transactions) == 0 {
		return components.MutedText.Render("No transactions")
	}

	visibleRows := components.Max(1, (height+1)/transactionRowHeight)
	m.ensureSelectedVisible(visibleRows)

	end := components.Min(len(m.transactions), m.offset+visibleRows)
	rows := make([]string, 0, end-m.offset)
	for i := m.offset; i < end; i++ {
		item := m.transactions[i]
		item.SetSelected(i == m.selected)
		rows = append(rows, item.ViewWidth(width))
	}

	return strings.Join(rows, "\n\n")
}

func (m *Model) View() string {
	return m.ViewWidthHeight(80, 20)
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
