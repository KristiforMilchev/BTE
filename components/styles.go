package components

import "github.com/charmbracelet/lipgloss"

var (
	Border          = lipgloss.Color("#263241")
	BorderOverlayer = lipgloss.Color("#2b3746")
	Muted           = lipgloss.Color("#6F7D8C")
	Text            = lipgloss.Color("#DDE5EE")
	Accent          = lipgloss.Color("#4A90E2")
	Success         = lipgloss.Color("#3BA55D")
	Warning         = lipgloss.Color("#F0B429")
	Danger          = lipgloss.Color("#D64545")

	App = lipgloss.NewStyle().Foreground(Text)

	SectionTitle  = lipgloss.NewStyle().Foreground(Accent).Bold(true)
	Label         = lipgloss.NewStyle().Foreground(Muted)
	Value         = lipgloss.NewStyle().Foreground(Text).Bold(true)
	MutedText     = lipgloss.NewStyle().Foreground(Muted)
	HelpTextStyle = lipgloss.NewStyle().Foreground(Muted)
	AmountText    = lipgloss.NewStyle().Foreground(Text).Bold(true)
	HeroAmount    = lipgloss.NewStyle().Height(2).AlignHorizontal(lipgloss.Center).Foreground(Text).Bold(true)
	ErrorText     = lipgloss.NewStyle().Foreground(Danger).Bold(true)
	SuccessText   = lipgloss.NewStyle().Foreground(Success).Bold(true)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Border).
			Padding(1, 2)
)
