package components_test

import (
	"testing"

	tokenlist "bos/components/token_list"
)

func TestTokenListZeroValueSelectedAsset(t *testing.T) {
	var model tokenlist.Model
	got := model.SelectedAsset()
	if got.Symbol != "ETH" || got.Balance != "0" || !got.Native {
		t.Fatalf("zero-value SelectedAsset() = %+v, want native ETH placeholder", got)
	}
}
