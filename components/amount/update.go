package amount

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Msg, tea.Cmd) {
	current, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil, nil
	}

	if current.String() == "enter" {
		return m, nil
	}

	if !isValidAmountChar(current.String()) {
		return nil, nil
	}

	data, _ := m.amountInput.Update(msg)
	m.amountInput = data
	return nil, nil
}

func isValidAmountChar(s string) bool {
	if s == "backspace" {
		return true
	}

	if len(s) != 1 {
		return false
	}
	c := s[0]
	return (c >= '0' && c <= '9') || c == '.'
}
