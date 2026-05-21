package views

import (
	"strings"

	"bos/components"
)

func RenderContacts(s State, width int) string {
	if len(s.Contacts) == 0 {
		return components.MutedText.Render("No saved recipients")
	}

	rows := make([]string, 0, len(s.Contacts))
	for i, c := range s.Contacts {
		active := s.Focus == FocusContacts && i == s.SelectedContact
		selected := i == s.SelectedContact

		nameStyle := components.Value
		marker := "  "
		if selected {
			nameStyle = nameStyle.Copy().Foreground(components.Accent)
		}
		if active {
			marker = components.FocusMarker(true)
		}

		row := marker + nameStyle.Render(components.Truncate(c.Name, width-2)) + "\n" +
			"  " + components.MutedText.Render(components.ShortAddress(c.Address))
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n\n")
}
