package networkDialog

import (
	"bos/components"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if !m.Visible {
		return ""
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.field("Name", 0),
		m.field("RPC", 1),
		m.field("Symbol", 2),
		m.field("Chain Id", 3),
		m.field("Block Explorer", 4),
		"",
		helpStyle.Render("Tab/Enter next • Shift+Tab back • Esc cancel"),
	)

	return dialogStyle.Render(content)
}

func (m *Model) field(label string, index int) string {
	if index < 0 || index >= len(m.inputs) {
		return ""
	}

	border := lipgloss.RoundedBorder()

	top := "─ " + label + " "

	width := 48
	remaining := width - lipgloss.Width(top) - 2
	if remaining < 0 {
		remaining = 0
	}

	border.Top = top + strings.Repeat("─", remaining)

	style := inputBoxStyle.Border(border)

	if m.focus == index {
		style = style.BorderForeground(components.Accent)
	}

	return style.Render(m.inputs[index].View())
}
