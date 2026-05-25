package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
func Truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= width {
		return s
	}
	if width <= 1 {
		return "…"
	}

	out := ""
	for _, r := range s {
		next := out + string(r)
		if lipgloss.Width(next+"…") > width {
			break
		}
		out = next
	}

	return out + "…"
}

func ShortAddress(address string) string {
	address = strings.TrimSpace(address)
	if len(address) <= 16 {
		return address
	}
	return address[:8] + "..." + address[len(address)-6:]
}

func Separator(width int) string {
	width = Max(1, width)

	return lipgloss.NewStyle().
		Width(width).
		MaxWidth(width).
		Foreground(Border).
		Render(strings.Repeat("─", width))
}
func KeyValue(key string, value string, width int) string {
	keyWidth := 14
	if width > 0 && width < 24 {
		keyWidth = 10
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		Label.Width(keyWidth).Render(key),
		Value.Render(value),
	)
}

func FocusMarker(active bool) string {
	return "  "
}

func Button(label string, active bool) string {
	if active {
		return lipgloss.NewStyle().Foreground(Accent).Bold(true).Render("[ " + label + " ]")
	}
	return MutedText.Render("[ " + label + " ]")
}

func Panel(width int, body string) string {
	return PanelSized(width, 0, body)
}

func PanelSized(width int, height int, body string) string {
	width = Max(12, width)

	innerWidth := Max(1, width-PanelStyle.GetHorizontalFrameSize())
	style := PanelStyle.Width(innerWidth)

	if height > 0 {
		innerHeight := Max(1, height-PanelStyle.GetVerticalFrameSize())
		style = style.Height(innerHeight)
	}

	return style.Render(body)
}

func HelpText(value string) string {
	return HelpTextStyle.Render(value)
}

func AmountBox(width int, value string) string {
	inner := Max(8, width-2)
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(Border).
		Width(inner).
		Padding(0, 1).
		Render(value)
}
