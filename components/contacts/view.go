package contacts

import (
	"bos/components"
	"strings"
)

func (m *Model) Render(width int) string {
	if len(m.contacts) == 0 {
		return components.MutedText.Render("No saved recipients")
	}

	rows := make([]string, 0, len(m.contacts))
	for i, contact := range m.contacts {
		active := i == m.selectedContact
		selected := i == m.selectedContact

		nameStyle := components.Value
		marker := "  "
		if selected {
			nameStyle = nameStyle.Foreground(components.Accent)
		}
		if active {
			marker = components.FocusMarker(true)
		}

		row := marker + nameStyle.Render(contact.Name) + "\n" +
			"  " + components.MutedText.Render(components.ShortAddress(contact.Address))
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n\n")
}

func (m *Model) View() string {
	return m.Render(80)
}

func (m *Model) ViewWidth(width int) string {
	return m.Render(width)
}
