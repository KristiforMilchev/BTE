package networksPopup

import (
	"bos/components/search"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

type SubmittedMsg struct {
	Network types.Network
}

type CancelledMsg struct{}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	if m.search == nil {
		m = New()
	}

	switch msg := msg.(type) {
	case search.SubmittedMsg:
		m.applySearch(msg.Query)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return CancelledMsg{} }
		case "f", "F":
			m.focusSearch()
			return m, nil
		case "tab":
			if m.focus == focusSearch {
				m.focusTable()
			} else {
				m.focusSearch()
			}
			return m, nil
		}
	}

	if m.focus == focusSearch {
		msg, cmd := m.search.Update(msg)
		m.filterNetworks(m.search.Value())
		if msg != nil {
			return m, func() tea.Msg { return msg }
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k", "up":
			m.moveSelection(-1)
			return m, nil
		case "j", "down":
			m.moveSelection(1)
			return m, nil
		case "enter", "space":
			return m, m.submit()
		}
	}

	return m, nil
}

func (m *Model) moveSelection(delta int) {
	if len(*m.filtered) == 0 {
		return
	}

	m.selected += delta
	if m.selected < 0 {
		m.selected = 0
	}
	if m.selected >= len(*m.filtered) {
		m.selected = len(*m.filtered) - 1
	}
}
