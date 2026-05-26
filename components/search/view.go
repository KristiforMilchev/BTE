package search

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) ViewWidth(width int, focused bool) string {
	width = components.Max(12, width)

	border := lipgloss.RoundedBorder()
	top := "─ Search "
	remaining := width - lipgloss.Width(top) - 2
	if remaining < 0 {
		remaining = 0
	}
	border.Top = top + strings.Repeat("─", remaining)

	style := lipgloss.NewStyle().
		Width(components.Max(1, width-2)).
		Padding(0, 1).
		Border(border).
		BorderForeground(components.Text)

	if focused {
		style = style.BorderForeground(components.Accent)
	}

	return style.Render(m.input.View())
}

func (m *Model) View() string {
	return m.ViewWidth(48, m.input.Focused())
}
