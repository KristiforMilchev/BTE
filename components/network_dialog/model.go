package networkDialog

import (
	"bos/components"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	dialogStyle = lipgloss.NewStyle().
			Width(52).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(components.BorderOverlayer)

	inputBoxStyle = lipgloss.NewStyle().
			Width(46).
			Padding(0, 1).
			MarginBottom(1).
			BorderForeground(components.Text)

	helpStyle = lipgloss.NewStyle().
			Foreground(components.Muted)
)

type Network struct {
	Name    string
	RPC     string
	Symbol  string
	ChainID int64
}

type SubmittedMsg struct {
	Network Network
}

type CancelledMsg struct{}

type Model struct {
	Visible bool
	inputs  []textinput.Model
	focus   int
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func New() *Model {
	labels := []string{"Name", "RPC", "Symbol", "Chain Id"}
	inputs := make([]textinput.Model, len(labels))

	for i, label := range labels {
		t := textinput.New()
		t.Placeholder = label
		t.CharLimit = 256
		t.Width = 42

		if i == 0 {
			t.Focus()
		}

		inputs[i] = t
	}

	return &Model{
		Visible: false,
		inputs:  inputs,
		focus:   0,
	}
}

func (m *Model) next() {
	if len(m.inputs) == 0 {
		return
	}

	m.inputs[m.focus].Blur()
	m.focus = (m.focus + 1) % len(m.inputs)
	m.inputs[m.focus].Focus()
}

func (m *Model) prev() {
	if len(m.inputs) == 0 {
		return
	}

	m.inputs[m.focus].Blur()
	m.focus--

	if m.focus < 0 {
		m.focus = len(m.inputs) - 1
	}

	m.inputs[m.focus].Focus()
}

func (m *Model) submit() tea.Cmd {
	return func() tea.Msg {
		chainID, err := strconv.ParseInt(strings.TrimSpace(m.inputs[3].Value()), 10, 64)
		if err != nil {
			chainID = 0
		}

		return SubmittedMsg{
			Network: Network{
				Name:    strings.TrimSpace(m.inputs[0].Value()),
				RPC:     strings.TrimSpace(m.inputs[1].Value()),
				Symbol:  strings.TrimSpace(m.inputs[2].Value()),
				ChainID: chainID,
			},
		}
	}
}
