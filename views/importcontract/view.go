package importcontract

import (
	"strings"

	"bos/components"
	"bos/enums"
	"bos/types"
	"bos/views"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	return views.RenderApp(m.width, m.height, enums.FocusTokens, m.status, func(width, height int) string {
		return m.render(width, height)
	})
}

func (m *Model) render(width int, height int) string {
	outerGap := 2
	panelGap := 2
	contentWidth := components.Max(72, width-(outerGap*2))
	contentHeight := components.Max(1, height)
	usableHeight := components.Max(1, contentHeight-1)
	panelWidth := (contentWidth - panelGap*2) / 3
	centerWidth := panelWidth
	m.input.Prompt = ""
	m.input.Width = components.Max(18, centerWidth-importPanelFrameWidth())

	addressPanel := importPanelSized(centerWidth, 6, strings.Join([]string{
		components.SectionTitle.
			Width(importPanelBodyWidth(centerWidth)).
			AlignHorizontal(lipgloss.Center).
			Render("Import Address"),
		components.Separator(importPanelBodyWidth(centerWidth)),
		m.input.View(),
		components.HelpText("enter load • S simulate • Y save • N cancel"),
	}, "\n"))

	sideHeight := usableHeight
	readableHeight := components.Max(8, usableHeight-lipgloss.Height(addressPanel)-1)
	centerColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		addressPanel,
		"",
		m.methodsPanel("Readable Methods", m.readable, centerWidth, readableHeight),
	)

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.methodsPanel("Callable Methods", m.callable, panelWidth, sideHeight),
		strings.Repeat(" ", panelGap),
		centerColumn,
		strings.Repeat(" ", panelGap),
		m.interactionsPanel(panelWidth, sideHeight),
	)

	return lipgloss.NewStyle().
		Width(width).
		Height(contentHeight).
		PaddingLeft(outerGap).
		PaddingRight(outerGap).
		Render(body)
}

func (m *Model) methodsPanel(title string, methods []Method, width int, height int) string {
	bodyWidth := importPanelBodyWidth(width)
	rows := []string{
		components.SectionTitle.
			Width(bodyWidth).
			AlignHorizontal(lipgloss.Center).
			Render(title),
		components.Separator(bodyWidth),
	}

	if !m.loaded {
		rows = append(rows, "", components.MutedText.Render("No contract address selected"))
	} else if len(methods) == 0 {
		rows = append(rows, "", components.MutedText.Render("No methods discovered"))
	} else {
		rows = append(rows, methodRows(methods, bodyWidth)...)
	}

	return importPanelSized(width, height, strings.Join(rows, "\n"))
}

func methodRows(methods []Method, width int) []string {
	rows := make([]string, 0, len(methods))
	for _, method := range methods {
		rows = append(rows, components.Value.Copy().Bold(false).Render(
			components.Truncate(method.Name, components.Max(1, width/3))+"  "+
				components.MutedText.Render(components.Truncate(method.Signature, components.Max(1, width-width/3-2))),
		))
	}
	return rows
}

func (m *Model) interactionsPanel(width int, height int) string {
	bodyWidth := importPanelBodyWidth(width)
	rows := []string{
		components.SectionTitle.
			Width(bodyWidth).
			AlignHorizontal(lipgloss.Center).
			Render("24h Interactions"),
		components.Separator(bodyWidth),
	}

	if !m.loaded {
		rows = append(rows, "", components.MutedText.Render("No contract address selected"))
	} else if m.loadingInteractions {
		rows = append(rows, "", components.MutedText.Render("Loading 24h interactions..."))
	} else if len(m.interactions) == 0 {
		rows = append(rows, "", components.MutedText.Render("No interactions loaded"))
	} else {
		rows = append(rows, interactionRows(m.interactions, bodyWidth)...)
	}

	return importPanelSized(width, height, strings.Join(rows, "\n"))
}

func importPanelBodyWidth(width int) int {
	return components.Max(1, width-importPanelFrameWidth()-importPanelPaddingWidth())
}

func importPanelFrameWidth() int {
	return importPanelStyle().GetHorizontalFrameSize()
}

func importPanelPaddingWidth() int {
	return importPanelStyle().GetHorizontalPadding()
}

func importPanelSized(width int, height int, body string) string {
	width = components.Max(12, width)
	style := importPanelStyle().Width(components.Max(1, width-importPanelStyle().GetHorizontalBorderSize()))
	if height > 0 {
		style = style.Height(components.Max(1, height-importPanelStyle().GetVerticalBorderSize()))
	}
	return style.Render(body)
}

func importPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(components.Border).
		Padding(0, 2)
}

func interactionRows(interactions []types.ContractInteraction, width int) []string {
	addressWidth := components.Max(12, width-14)
	rows := make([]string, 0, len(interactions)+2)
	rows = append(rows, components.Label.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(addressWidth).Render("Address"),
			lipgloss.NewStyle().Width(8).Render("Action"),
			lipgloss.NewStyle().Width(6).Render("Age"),
		),
	))
	for _, interaction := range interactions {
		rows = append(rows, components.Value.Copy().Bold(false).Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.NewStyle().Width(addressWidth).Render(components.TruncateMiddle(interaction.Address, addressWidth-1)),
				lipgloss.NewStyle().Width(8).Render(components.Truncate(interaction.Action, 8)),
				lipgloss.NewStyle().Width(6).Render(components.Truncate(interaction.Age, 6)),
			),
		))
	}
	return rows
}
