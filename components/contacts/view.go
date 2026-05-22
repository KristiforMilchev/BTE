package contacts

import (
	"bos/components"
	"bos/layout"
	"strings"
)

func (m *Model) Render(ctx layout.Context) string {
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
	return m.Render(layout.Context{Constraints: layout.Constraints{Width: 80, Height: 24}})
}
