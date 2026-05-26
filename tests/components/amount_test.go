package components_test

import (
	"strings"
	"testing"

	"bos/components/amount"
	"bos/types"
)

func TestAmountUpdateAcceptsOnlyAmountCharacters(t *testing.T) {
	model := amount.New()
	model.Focus()

	for _, input := range []string{"1", ".", "2"} {
		if msg, _ := model.Update(key(input)); msg != nil {
			t.Fatalf("Update(%q) returned message %T, want nil", input, msg)
		}
	}
	if got := model.Value(); got != "1.2" {
		t.Fatalf("Value() = %q, want %q", got, "1.2")
	}

	if msg, _ := model.Update(key("x")); msg != nil {
		t.Fatalf("Update(invalid key) returned message %T, want nil", msg)
	}
	if got := model.Value(); got != "1.2" {
		t.Fatalf("invalid key changed value to %q", got)
	}
}

func TestAmountUpdateEnterReturnsModelMessage(t *testing.T) {
	model := amount.New()

	msg, _ := model.Update(key("enter"))
	if msg != model {
		t.Fatalf("Update(enter) = %T, want current model", msg)
	}
}

func TestAmountClearResetsValue(t *testing.T) {
	model := amount.New()
	model.Focus()

	for _, input := range []string{"1", ".", "2"} {
		if msg, _ := model.Update(key(input)); msg != nil {
			t.Fatalf("Update(%q) returned message %T, want nil", input, msg)
		}
	}

	model.Clear()

	if got := model.Value(); got != "" {
		t.Fatalf("Value() after Clear() = %q, want empty", got)
	}
}

func TestAmountViewDoesNotRenderTrailingBlankLine(t *testing.T) {
	model := amount.New()
	model.SetSymbol(types.Token{Symbol: "ETH"})
	view := model.View()
	lines := strings.Split(view, "\n")

	if strings.TrimSpace(lines[len(lines)-1]) == "" {
		t.Fatal("amount view ended with a blank line")
	}
}
