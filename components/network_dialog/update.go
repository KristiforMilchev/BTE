package networkDialog

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	if len(m.inputs) == 0 {
		m = New()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return CancelledMsg{}
			}

		case "tab", "down":
			m.next()
			return m, nil

		case "shift+tab", "up":
			m.prev()
			return m, nil

		case "enter", "space":
			if m.focus == len(m.inputs)-1 {
				return m, m.submit()
			}

			m.next()
			return m, nil
		}
	}

	var cmds []tea.Cmd
	for i := range m.inputs {
		if i != m.focus {
			continue
		}

		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
