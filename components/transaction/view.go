package transaction

import (
	"strings"

	"bos/components"
)

func (m Model) ViewWidth(width int) string {
	width = components.Max(12, width)

	marker := "  "
	valueStyle := components.Value
	if m.selected {
		marker = components.FocusMarker(true)
		valueStyle = valueStyle.Copy().Foreground(components.Accent)
	}

	valueWidth := components.Max(1, width-2)
	tx := m.transaction

	return strings.Join([]string{
		marker + valueStyle.Render("To: "+components.Truncate(tx.To, valueWidth)),
		"  " + components.MutedText.Render("Block: "+components.Truncate(tx.Block, valueWidth)),
		"  " + components.MutedText.Render("TxHash: "+components.Truncate(tx.TxHash, valueWidth)),
		"  " + components.MutedText.Render("Amount: "+components.Truncate(tx.Amount, valueWidth)),
	}, "\n")
}
