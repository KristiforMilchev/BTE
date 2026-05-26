package networksPopup

import (
	"strconv"
	"strings"

	"bos/components"
	"bos/types"

	"github.com/charmbracelet/lipgloss"
)

var (
	dialogStyle = lipgloss.NewStyle().
			Width(76).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(components.BorderOverlayer)

	helpStyle = lipgloss.NewStyle().
			Foreground(components.Muted)
)

const (
	tableHeight  = 10
	nameWidth    = 16
	rpcWidth     = 38
	chainIDWidth = 10
)

func (m *Model) View() string {
	if !m.Visible {
		return ""
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.search.ViewWidth(70, m.focus == focusSearch),
		"",
		m.tableView(tableHeight),
		"",
		helpStyle.Render("F search • j/k move • enter/space select • tab switch • esc cancel"),
	)

	return dialogStyle.Render(content)
}

func (m *Model) tableView(height int) string {
	header := components.Value.Render(
		padCell("NAME", nameWidth) + "  " +
			padCell("RPC", rpcWidth) + "  " +
			padCell("CHAINID", chainIDWidth),
	)
	separator := components.Separator(nameWidth + rpcWidth + chainIDWidth + 4)

	rows := []string{header, separator}
	if m.filtered == nil || len(*m.filtered) == 0 {
		rows = append(rows, components.MutedText.Render("No networks found"))
		return strings.Join(rows, "\n")
	}

	visibleRows := components.Max(1, height-2)
	m.ensureSelectedVisible(visibleRows)

	end := components.Min(len(*m.filtered), m.offset+visibleRows)
	for i := m.offset; i < end; i++ {
		network := (*m.filtered)[i]
		row := formatNetworkRow(&network)
		if i == m.selected {
			row = components.FocusMarker(true) + lipgloss.NewStyle().
				Foreground(components.Accent).
				Bold(true).
				Render(row)
		} else {
			row = "  " + row
		}
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n")
}

func (m *Model) ensureSelectedVisible(visibleRows int) {
	if visibleRows < 1 {
		visibleRows = 1
	}
	if m.selected < m.offset {
		m.offset = m.selected
	}
	if m.selected >= m.offset+visibleRows {
		m.offset = m.selected - visibleRows + 1
	}
	if maxOffset := components.Max(0, len(*m.filtered)-visibleRows); m.offset > maxOffset {
		m.offset = maxOffset
	}
}

func formatNetworkRow(network *types.Network) string {
	chain := strconv.FormatInt(network.Chain.Int64(), 10)
	return padCell(*network.Name, nameWidth) + "  " +
		padCell(*network.Rpc, rpcWidth) + "  " +
		padCell(chain, chainIDWidth)
}

func padCell(value string, width int) string {
	value = components.Truncate(value, width)
	padding := width - lipgloss.Width(value)
	if padding < 0 {
		padding = 0
	}
	return value + strings.Repeat(" ", padding)
}
