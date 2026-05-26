package transactions

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k", "up":
			m.moveSelection(-1)
			return nil, nil
		case "j", "down":
			m.moveSelection(1)
			return nil, nil
		case "enter", "space":
			return m, nil
		}
	}

	return nil, nil
}

func (m *Model) moveSelection(delta int) {
	if len(m.transactions) == 0 {
		return
	}

	m.selected += delta
	if m.selected < 0 {
		m.selected = 0
	}
	if m.selected >= len(m.transactions) {
		m.selected = len(m.transactions) - 1
	}
}
