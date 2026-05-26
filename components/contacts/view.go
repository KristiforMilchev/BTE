package contacts

import (
	"bos/components"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Render(width int) string {
	width = components.Max(24, width)
	if len(m.contacts) == 0 {
		return components.MutedText.Render("No saved recipients")
	}

	rows := make([]string, 0, len(m.contacts)+2)
	rows = append(rows, contactHeader(width), components.Separator(width))
	for i, contact := range m.contacts {
		rows = append(rows, contactRow(contact.Name, contact.Address, width, i == m.selectedContact))
	}

	return strings.Join(rows, "\n")
}

func (m *Model) View() string {
	return m.Render(80)
}

func (m *Model) ViewWidth(width int) string {
	return m.Render(width)
}

func contactHeader(width int) string {
	nameWidth, addressWidth := contactColumnWidths(width)
	return components.Label.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(nameWidth).Render("Name"),
			lipgloss.NewStyle().Width(addressWidth).Render("Address"),
		),
	)
}

func contactRow(name string, address string, width int, selected bool) string {
	nameWidth, addressWidth := contactColumnWidths(width)
	style := components.Value.Copy().Bold(false)
	if selected {
		style = style.Foreground(components.Accent).Bold(true)
	}

	return style.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(nameWidth).Render(components.Truncate(name, nameWidth-1)),
			lipgloss.NewStyle().Width(addressWidth).Render(components.Truncate(address, addressWidth)),
		),
	)
}

func contactColumnWidths(width int) (int, int) {
	width = components.Max(32, width)
	addressWidth := components.Min(42, components.Max(8, width/2))
	nameWidth := components.Max(8, width-addressWidth)
	return nameWidth, addressWidth
}
