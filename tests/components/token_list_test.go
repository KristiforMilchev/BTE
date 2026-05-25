package components_test

import (
	"strings"
	"testing"

	tokenlist "bos/components/token_list"
	"bos/types"
)

func TestTokenListZeroValueSelectedAsset(t *testing.T) {
	var model tokenlist.Model
	got := model.SelectedAsset()
	if got.Symbol != "ETH" || got.Balance != "0" || !got.Native {
		t.Fatalf("zero-value SelectedAsset() = %+v, want native ETH placeholder", got)
	}
}

func TestTokenListSelectionScrollsIntoView(t *testing.T) {
	model := &tokenlist.Model{}
	model.SetTokens([]types.Token{
		{Symbol: "ETH", Balance: "1", Native: true},
		{Symbol: "USDC", Balance: "2", Verified: true},
		{Symbol: "DAI", Balance: "3", Verified: true},
		{Symbol: "WBTC", Balance: "4", Verified: true},
	})

	for range 3 {
		if _, cmd := model.Update(key("down")); cmd != nil {
			t.Fatal("down returned command, want nil")
		}
	}

	view := model.ViewWidthHeight(42, 4)
	if !strings.Contains(view, "WBTC") {
		t.Fatalf("scrolled view = %q, want selected token visible", view)
	}
	if strings.Contains(view, "ETH") {
		t.Fatalf("scrolled view = %q, want first token outside viewport", view)
	}
}
