package views

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func selectedAsset(s State) Token {
	if len(s.Tokens) == 0 {
		return Token{Symbol: "ETH", Balance: "0", Native: true}
	}
	if s.SelectedToken < 0 || s.SelectedToken >= len(s.Tokens) {
		return s.Tokens[0]
	}
	return s.Tokens[s.SelectedToken]
}

func selectedRecipient(s State) Contact {
	if len(s.Contacts) == 0 {
		return Contact{Name: "No Contact", Address: ""}
	}
	if s.SelectedContact < 0 || s.SelectedContact >= len(s.Contacts) {
		return s.Contacts[0]
	}
	return s.Contacts[s.SelectedContact]
}

func riskLabel(risk string) string {
	switch strings.ToLower(risk) {
	case "low":
		return lipgloss.NewStyle().Foreground(components.Success).Bold(true).Render(risk)
	case "medium":
		return lipgloss.NewStyle().Foreground(components.Warning).Bold(true).Render(risk)
	case "high", "critical":
		return lipgloss.NewStyle().Foreground(components.Danger).Bold(true).Render(risk)
	default:
		return risk
	}
}
