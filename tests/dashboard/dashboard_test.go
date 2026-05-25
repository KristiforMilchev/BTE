package dashboard_test

import (
	"strings"
	"testing"

	"bos/components"
	"bos/views/dashboard"

	"github.com/charmbracelet/lipgloss"
)

func TestPanelContentSizedUsesRequestedOuterSize(t *testing.T) {
	const width = 32
	const height = 8

	rendered := dashboard.PanelContentSized(width, height, "body")
	lines := strings.Split(rendered, "\n")

	if got := lipgloss.Height(rendered); got != height {
		t.Fatalf("PanelContentSized height = %d, want %d", got, height)
	}

	for i, line := range lines {
		if got := lipgloss.Width(line); got != width {
			t.Fatalf("line %d width = %d, want %d", i, got, width)
		}
	}
}

func TestPanelBodyWidthMatchesPanelContentArea(t *testing.T) {
	const width = 32
	want := width - components.PanelStyle.GetHorizontalFrameSize()

	if got := dashboard.PanelBodyWidth(width); got != want {
		t.Fatalf("PanelBodyWidth(%d) = %d, want %d", width, got, want)
	}
}
