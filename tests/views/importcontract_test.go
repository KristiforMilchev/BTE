package views_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"bos/types"
	"bos/views/importcontract"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethereum/go-ethereum/common"
)

type mockContractInteractionReader struct {
	called bool
}

func (m *mockContractInteractionReader) RecentInteractions(_ context.Context, contract common.Address, since time.Time) ([]types.ContractInteraction, error) {
	m.called = true
	if contract != common.HexToAddress("0x1111111111111111111111111111111111111111") {
		return nil, nil
	}
	if time.Since(since) < 23*time.Hour {
		return nil, nil
	}
	return []types.ContractInteraction{
		{
			Address: "0x2222222222222222222222222222222222222222",
			Action:  "Transfer",
			Age:     "2h",
			TxHash:  "0xabc",
		},
	}, nil
}

func TestImportContractViewShowsEmptyState(t *testing.T) {
	model := importcontract.New()
	next, _ := model.Update(tea.WindowSizeMsg{Width: 140, Height: 40})

	rendered := next.View()
	if got := lipgloss.Height(rendered); got != 40 {
		t.Fatalf("rendered height = %d, want 40", got)
	}
	if !strings.Contains(rendered, "Import Address") {
		t.Fatalf("rendered view missing title: %q", rendered)
	}
	if !strings.Contains(rendered, "No contract address selected") {
		t.Fatalf("rendered view missing empty state: %q", rendered)
	}
}

func TestImportContractViewLoadsPlaceholderSections(t *testing.T) {
	reader := &mockContractInteractionReader{}
	model := importcontract.New(reader)
	next, _ := model.Update(tea.WindowSizeMsg{Width: 140, Height: 40})
	model = next.(*importcontract.Model)

	for _, r := range "0x1111111111111111111111111111111111111111" {
		next, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = next.(*importcontract.Model)
	}
	next, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("enter returned nil command, want interaction load command")
	}
	if rendered := next.View(); !strings.Contains(rendered, "Loading 24h interactions") {
		t.Fatalf("rendered view missing loading state: %q", rendered)
	}
	model = next.(*importcontract.Model)
	next, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Z")})
	if strings.Contains(next.View(), "Z") {
		t.Fatalf("input accepted typing after enter blur: %q", next.View())
	}
	model = next.(*importcontract.Model)
	next, _ = model.Update(cmd())

	rendered := next.View()
	for _, want := range []string{"Callable Methods", "Readable Methods", "24h Interactions", "transfer", "balanceOf", "Transfer", "0x2222"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered view missing %q: %q", want, rendered)
		}
	}
	if !reader.called {
		t.Fatal("interaction reader was not called")
	}
}

func TestImportContractClearRefocusesInput(t *testing.T) {
	model := importcontract.New()
	next, _ := model.Update(tea.WindowSizeMsg{Width: 140, Height: 40})
	model = next.(*importcontract.Model)

	for _, r := range "0x1111111111111111111111111111111111111111" {
		next, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = next.(*importcontract.Model)
	}
	next, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = next.(*importcontract.Model)

	next, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Z")})
	if strings.Contains(next.View(), "Z") {
		t.Fatalf("input accepted typing while blurred: %q", next.View())
	}
	model = next.(*importcontract.Model)

	next, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("C")})
	model = next.(*importcontract.Model)
	next, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	if !strings.Contains(next.View(), "a") {
		t.Fatalf("input did not accept typing after clear/refocus: %q", next.View())
	}
}
