package contacts

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "k", "up":
			if m.selectedContact <= 0 {
				m.selectedContact = 0
				return nil, nil
			}

			m.selectedContact--
		case "j", "down":
			if m.selectedContact > len(m.contacts) {
				m.selectedContact = len(m.contacts) - 1
				return nil, nil
			}

			m.selectedContact++
		case "enter", "space":

			return m, nil
		}
	}
	return nil, nil

}
