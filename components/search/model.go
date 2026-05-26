package search

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SubmittedMsg struct {
	Query string
}

type Model struct {
	input textinput.Model
}

func New(placeholder string, width int) *Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = 256
	input.Width = width
	input.Focus()

	return &Model{input: input}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Focus() {
	m.input.Focus()
}

func (m *Model) Blur() {
	m.input.Blur()
}

func (m *Model) Value() string {
	return m.input.Value()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			return nil, func() tea.Msg {
				return SubmittedMsg{Query: strings.TrimSpace(m.input.Value())}
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return nil, cmd
}
