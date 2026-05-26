package components_test

import (
	"context"
	"testing"

	contactDialog "bos/components/contact_dialog"
	"bos/di"
	"bos/repositories"
	"bos/tests/testmocks"

	tea "github.com/charmbracelet/bubbletea"
)

func TestContactDialogCancelCommand(t *testing.T) {
	model := contactDialog.New()
	_, cmd := model.Update(key("esc"))
	if cmd == nil {
		t.Fatal("Update(esc) returned nil command, want cancellation command")
	}
	if _, ok := cmd().(contactDialog.CancelledMsg); !ok {
		t.Fatalf("cancel command returned %T, want contactDialog.CancelledMsg", cmd())
	}
}

func TestContactDialogSubmissionFromTypedInput(t *testing.T) {
	storage := testmocks.NewStorage(t)
	register := repositories.NewRegister(storage)
	di.SetupDependenciesWith(di.Dependencies{Storage: storage, Register: &register})
	t.Cleanup(func() {
		di.SetupDependenciesWith(di.Dependencies{})
	})

	model := contactDialog.New()

	for _, input := range []tea.KeyMsg{
		key("A"), key("l"), key("i"), key("c"), key("e"),
		key("tab"),
		key("0"), key("x"), key("1"), key("1"), key("1"),
	} {
		var cmd tea.Cmd
		model, cmd = model.Update(input)
		if cmd != nil {
			_ = cmd()
		}
	}

	_, cmd := model.Update(key("enter"))
	if cmd == nil {
		t.Fatal("Update(enter) returned nil command, want submission command")
	}
	msg, ok := cmd().(contactDialog.SubmittedMsg)
	if !ok {
		t.Fatalf("submit command returned %T, want contactDialog.SubmittedMsg", msg)
	}
	if msg.Contact.Name != "Alice" || msg.Contact.Address != "0x111" {
		t.Fatalf("submitted contact = %+v, want typed values", msg.Contact)
	}

	var name string
	err := storage.QueryRow(context.Background(), "SELECT name FROM contacts WHERE address = ?;", "0x111").Scan(&name)
	if err != nil {
		t.Fatalf("saved contact lookup returned error: %v", err)
	}
	if name != "Alice" {
		t.Fatalf("saved contact name = %q, want Alice", name)
	}
}
