package components_test

import (
	"testing"

	networkDialog "bos/components/network_dialog"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNetworkDialogCancelCommand(t *testing.T) {
	model := networkDialog.New()
	_, cmd := model.Update(key("esc"))
	if cmd == nil {
		t.Fatal("Update(esc) returned nil command, want cancellation command")
	}
	if _, ok := cmd().(networkDialog.CancelledMsg); !ok {
		t.Fatalf("cancel command returned %T, want networkDialog.CancelledMsg", cmd())
	}
}

func TestNetworkDialogSubmissionFromTypedInput(t *testing.T) {
	model := networkDialog.New()

	for _, input := range []tea.KeyMsg{
		key("L"), key("o"), key("c"), key("a"), key("l"),
		key("tab"),
		key("h"), key("t"), key("t"), key("p"),
		key("tab"),
		key("E"), key("T"), key("H"),
		key("tab"),
		key("1"), key("3"), key("3"), key("7"),
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
	msg, ok := cmd().(networkDialog.SubmittedMsg)
	if !ok {
		t.Fatalf("submit command returned %T, want networkDialog.SubmittedMsg", msg)
	}

	if msg.Network.Name != "Local" || msg.Network.RPC != "http" || msg.Network.Symbol != "ETH" || msg.Network.ChainID != 1337 {
		t.Fatalf("submitted network = %+v, want typed values", msg.Network)
	}
}
