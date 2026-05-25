package components_test

import (
	"strings"
	"testing"

	networksPopup "bos/components/network_popup"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNetworkPopupSearchFiltersAndSelectionSubmitsNetwork(t *testing.T) {
	model := networksPopup.New()
	model.Visible = true

	for _, input := range []tea.KeyMsg{key("b"), key("a"), key("s"), key("e")} {
		var cmd tea.Cmd
		model, cmd = model.Update(input)
		if cmd != nil {
			_ = cmd()
		}
	}

	liveView := model.View()
	if !strings.Contains(liveView, "Base") {
		t.Fatalf("live filtered view = %q, want Base network", liveView)
	}
	if strings.Contains(liveView, "Ethereum") {
		t.Fatalf("live filtered view = %q, want Ethereum filtered out", liveView)
	}

	var cmd tea.Cmd
	model, cmd = model.Update(key("enter"))
	if cmd == nil {
		t.Fatal("search enter returned nil command, want search submission")
	}
	msg := cmd()

	model, cmd = model.Update(msg)
	if cmd != nil {
		_ = cmd()
	}

	view := model.View()
	if !strings.Contains(view, "Base") {
		t.Fatalf("filtered view = %q, want Base network", view)
	}
	if strings.Contains(view, "Ethereum") {
		t.Fatalf("filtered view = %q, want Ethereum filtered out", view)
	}

	model, cmd = model.Update(key("space"))
	if cmd == nil {
		t.Fatal("space returned nil command, want network submission")
	}
	submitted, ok := cmd().(networksPopup.SubmittedMsg)
	if !ok {
		t.Fatalf("network submit returned %T, want networksPopup.SubmittedMsg", submitted)
	}
	if submitted.Network.Name != "Base" || submitted.Network.ChainID != 8453 {
		t.Fatalf("submitted network = %+v, want Base chain 8453", submitted.Network)
	}
}
