package main

import (
	"bos/di"
	"bos/views/app"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	di.SetupDependencies()
	defer di.GetLogger().Close()

	program := tea.NewProgram(
		app.New(di.GetWallet(), di.GetNetwork(), di.Repositories(), di.GetContractInteractions()),
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
