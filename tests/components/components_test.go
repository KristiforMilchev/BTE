package components_test

import (
	"strings"
	"testing"

	"bos/components/amount"
	"bos/components/contacts"
	networkDialog "bos/components/network_dialog"
	tokenlist "bos/components/token_list"
	transactionsComponent "bos/components/transactions"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

func key(s string) tea.KeyMsg {
	if len(s) == 1 {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}

	switch s {
	case "backspace":
		return tea.KeyMsg{Type: tea.KeyBackspace}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

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

func TestAmountViewDoesNotRenderTrailingBlankLine(t *testing.T) {
	model := amount.New()
	model.SetSymbol(types.Token{Symbol: "ETH"})
	view := model.View()
	lines := strings.Split(view, "\n")

	if strings.TrimSpace(lines[len(lines)-1]) == "" {
		t.Fatal("amount view ended with a blank line")
	}
}

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

func TestTokenListZeroValueSelectedAsset(t *testing.T) {
	var model tokenlist.Model
	got := model.SelectedAsset()
	if got.Symbol != "ETH" || got.Balance != "0" || !got.Native {
		t.Fatalf("zero-value SelectedAsset() = %+v, want native ETH placeholder", got)
	}
}

func TestTransactionsSelectionAndScroll(t *testing.T) {
	model := transactionsComponent.New([]types.Transaction{
		{To: "0x1111111111111111111111111111111111111111", Block: "1", TxHash: "0xaaa", Amount: "1 ETH"},
		{To: "0x2222222222222222222222222222222222222222", Block: "2", TxHash: "0xbbb", Amount: "2 ETH"},
		{To: "0x3333333333333333333333333333333333333333", Block: "3", TxHash: "0xccc", Amount: "3 ETH"},
	})

	if _, cmd := model.Update(key("down")); cmd != nil {
		t.Fatal("down returned command, want nil")
	}
	if _, cmd := model.Update(key("down")); cmd != nil {
		t.Fatal("second down returned command, want nil")
	}

	view := model.ViewWidthHeight(40, 5)
	if !strings.Contains(view, "3 ETH") {
		t.Fatalf("scrolled view = %q, want selected transaction visible", view)
	}
	if strings.Contains(view, "1 ETH") {
		t.Fatalf("scrolled view = %q, want first transaction outside viewport", view)
	}

	msg, _ := model.Update(key("space"))
	if msg == nil {
		t.Fatal("space returned nil, want selection message")
	}
}

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
