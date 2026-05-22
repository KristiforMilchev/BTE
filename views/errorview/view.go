package errorview

import (
	"bos/components"
	"bos/utils"
)

func (m *Model) View() string {
	title := m.payload.Title
	if title == "" {
		title = "Blockcert"
	}

	body := components.ErrorText.Render("Error:") + "\n" + m.payload.Message + "\n\nPress enter to retry. Press esc to return. Press q to quit."
	return utils.RenderCentered(m.width, m.height, title, body)
}
