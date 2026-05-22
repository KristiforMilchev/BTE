package views

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func RenderApp(width int, height int, focus FocusArea, statusMessage string, rpcURL string, content string) string {
	width = safeWidth(width)
	if height <= 0 {
		height = 30
	}

	return components.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			RenderHeader(width, focus),
			content,
			RenderFooter(width, statusMessage, rpcURL),
		),
	)
}

func RenderCentered(width int, height int, title string, body string) string {
	width = safeWidth(width)
	if height <= 0 {
		height = 30
	}

	content := components.SectionTitle.Render(title) + "\n\n" + body
	box := components.Panel(64, content)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func RenderHeader(width int, focus FocusArea) string {
	width = safeWidth(width)

	left := components.Value.Render("Blockcert")
	right := components.MutedText.Render(activeHelp(focus))

	spacerWidth := components.Max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-2)
	line := left + lipgloss.NewStyle().Width(spacerWidth).Render("") + right

	return line + "\n" + components.Separator(width)
}

func activeHelp(focus FocusArea) string {
	switch focus {
	case FocusAmount:
		return "Amount • type value • h/l assets • p recipients • s simulate • S send"
	case FocusTokens:
		return "Assets • j/k choose token • h amount • p recipients • s simulate"
	case FocusContacts:
		return "Recipients • j/k choose • enter select • h amount • l assets"
	default:
		return "hjkl move • p recipients • s simulate • S send"
	}
}

func RenderFooter(width int, statusMessage string, rpcURL string) string {
	width = safeWidth(width)

	status := "Ledger Connected • RPC Online"
	if rpcURL != "" {
		status += " • " + rpcURL
	}
	if strings.TrimSpace(statusMessage) != "" {
		status += " • " + statusMessage
	}

	return components.Separator(width) + "\n" + components.MutedText.Render(components.Truncate(status, width))
}

func safeWidth(width int) int {
	if width < 100 {
		return 100
	}
	return width
}
