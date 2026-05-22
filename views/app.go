package views

import (
	"strings"

	"bos/components"

	"github.com/charmbracelet/lipgloss"
)

func Render(s State) string {
	s.Width = safeWidth(s.Width)
	if s.Height <= 0 {
		s.Height = 30
	}

	switch s.Screen {
	case ScreenLoading:
		return RenderCentered(s, "Connecting to Ledger", "Requirements:\n- Ledger plugged in\n- Device unlocked\n- Ethereum app open\n- Ledger Live closed\n\nPress q to quit.")
	case ScreenDashboard:
		return renderApp(s, RenderDashboard(s))
	case ScreenConfirm:
		return renderApp(s, RenderConfirm(s))
	case ScreenSending:
		return renderApp(s, RenderSending(s))
	case ScreenSent:
		return renderApp(s, RenderSent(s))
	case ScreenError:
		return RenderCentered(s, "Blockcert", components.ErrorText.Render("Error:")+"\n"+s.Err+"\n\nPress enter to retry. Press esc to return. Press q to quit.")
	default:
		return ""
	}
}

func renderApp(s State, content string) string {
	return components.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			RenderHeader(s),
			content,
			RenderFooter(s),
		),
	)
}

func RenderCentered(s State, title string, body string) string {
	content := components.SectionTitle.Render(title) + "\n\n" + body
	box := components.Panel(64, content)

	return lipgloss.Place(
		s.Width,
		s.Height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func RenderHeader(s State) string {
	left := components.Value.Render("Blockcert")
	right := components.MutedText.Render(activeHelp(s))

	spacerWidth := components.Max(1, s.Width-lipgloss.Width(left)-lipgloss.Width(right)-2)
	line := left + lipgloss.NewStyle().Width(spacerWidth).Render("") + right

	return line + "\n" + components.Separator(s.Width)
}

func activeHelp(s State) string {
	switch s.Focus {
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

func RenderFooter(s State) string {
	status := "Ledger Connected • RPC Online"
	if s.RpcURL != "" {
		status += " • " + s.RpcURL
	}
	if strings.TrimSpace(s.StatusMessage) != "" {
		status += " • " + s.StatusMessage
	}

	return components.Separator(s.Width) + "\n" + components.MutedText.Render(components.Truncate(status, s.Width))
}

func safeWidth(width int) int {
	if width < 100 {
		return 100
	}
	return width
}
