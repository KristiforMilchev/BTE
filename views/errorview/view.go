package errorview

import (
	"bos/components"
	"bos/views"
)

func (m *Model) View() string {
	title := m.payload.Title
	if title == "" {
		title = "Blockcert"
	}

	body := components.ErrorText.Render("Error:") + "\n" + m.payload.Message + "\n\nPress enter to retry. Press esc to return. Press q to quit."
	return views.RenderCentered(m.width, m.height, title, body)
}
