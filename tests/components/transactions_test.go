package components_test

import (
	"strings"
	"testing"

	transactionsComponent "bos/components/transactions"
	"bos/types"
)

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
