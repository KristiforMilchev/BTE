package components_test

import (
	"testing"

	"bos/components/contacts"
)

func TestContactsSelection(t *testing.T) {
	model := contacts.NewContacts()
	if got := model.SelectedRecipient().Name; got != "Treasury Wallet" {
		t.Fatalf("initial recipient = %q, want %q", got, "Treasury Wallet")
	}

	if msg, _ := model.Update(key("down")); msg != nil {
		t.Fatalf("down returned message %T, want nil", msg)
	}
	if got := model.SelectedRecipient().Name; got != "Personal Wallet" {
		t.Fatalf("recipient after down = %q, want %q", got, "Personal Wallet")
	}

	msg, _ := model.Update(key("enter"))
	if msg == nil {
		t.Fatal("enter returned nil, want selection message")
	}
}

func TestEmptyContactsSelectionReturnsPlaceholder(t *testing.T) {
	model := &contacts.Model{}
	got := model.SelectedRecipient()
	if got.Name != "No Contact" || got.Address != "" {
		t.Fatalf("empty SelectedRecipient() = %+v, want no-contact placeholder", got)
	}
}
