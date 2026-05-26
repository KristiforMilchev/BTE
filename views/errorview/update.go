package errorview

import (
	"bos/enums"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, func() tea.Msg { return types.NavigateMsg{Screen: m.payload.Return} }
		case "esc":
			return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenDashboard} }
		}
	}
	return m, nil
}
