package contactDialog

import (
	"bos/components"
	"bos/di"
	"bos/types"
	"log"
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

type SubmittedMsg struct {
	Contact types.Contact
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
	labels := []string{"Name", "Address"}
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
		contact := types.Contact{
			Name:    strings.TrimSpace(m.inputs[0].Value()),
			Address: strings.TrimSpace(m.inputs[1].Value()),
		}

		if err := di.Repositories().Contacts.Create(contact); err != nil {
			log.Printf("Failed to save contact from dialog -> %s", err)
			return nil
		}

		return SubmittedMsg{Contact: contact}
	}
}
