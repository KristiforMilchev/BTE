package main

import (
	"bos/di"
	"bos/views"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	di.SetupDependencies()
	defer di.GetLogger().Close()

	amount := textinput.New()
	amount.Placeholder = "0.01"
	amount.CharLimit = 32
	amount.Width = 0
	amount.Focus()

	p := tea.NewProgram(&Model{
		screen:           views.ScreenLoading,
		focus:            views.FocusAmount,
		amountInput:      amount,
		selectedToken:    0,
		selectedContact:  0,
		tokens:           fakeTokens(),
		contacts:         fakeContacts(),
		simulationStatus: "Not Run",
		riskLevel:        "—",
		estimatedFee:     "—",
		statusMessage:    "Ready",
		network:          di.GetNetwork(),
		wallet:           di.GetWallet(),
	}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
