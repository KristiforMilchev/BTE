package components_test

import (
	"testing"

	"bos/components/search"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSearchSubmitReturnsQuery(t *testing.T) {
	model := search.New("Search", 20)

	for _, input := range []tea.KeyMsg{key("b"), key("a"), key("s"), key("e")} {
		if _, cmd := model.Update(input); cmd != nil {
			_ = cmd()
		}
	}

	_, cmd := model.Update(key("enter"))
	if cmd == nil {
		t.Fatal("enter returned nil command, want search submission")
	}
	msg, ok := cmd().(search.SubmittedMsg)
	if !ok {
		t.Fatalf("search submit returned %T, want search.SubmittedMsg", msg)
	}
	if msg.Query != "base" {
		t.Fatalf("search query = %q, want base", msg.Query)
	}
}
