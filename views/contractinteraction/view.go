package contractinteraction

import (
	"strings"

	"bos/components"
	"bos/enums"
	"bos/types"
	"bos/views"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	width := components.Clamp(m.width-8, 72, 112)
	contentWidth := components.Max(1, width-components.PanelStyle.GetHorizontalFrameSize())

	content := []string{
		components.SectionTitle.Render("Contract Interaction"),
		components.KeyValue("Address", components.ShortAddress(m.address), contentWidth),
		components.MutedText.Render("Fake ABI import for business-flow wiring."),
		"",
		components.SectionTitle.Render("Functions"),
		components.Separator(contentWidth),
		m.functionList(contentWidth),
		"",
		components.SectionTitle.Render("Parameters"),
		components.Separator(contentWidth),
		parameterList(m.functions[m.selected]),
		"",
		components.Button("enter/s  Simulate", true) + " " + components.Button("esc  Dashboard", false),
	}

	body := components.Panel(width, strings.Join(content, "\n"))
	return views.RenderApp(m.width, m.height, enums.FocusSend, "Contract interaction", func(width, height int) string {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, body)
	})
}

func (m *Model) functionList(width int) string {
	rows := make([]string, 0, len(m.functions))
	for i, function := range m.functions {
		style := components.Value
		if i == m.selected {
			style = style.Copy().Foreground(components.Accent)
		}
		row := function.Name + "  " + components.MutedText.Render(function.Signature)
		rows = append(rows, "  "+style.Render(components.Truncate(row, width-2)))
	}
	return strings.Join(rows, "\n")
}

func parameterList(function types.ContractFunction) string {
	if len(function.Parameters) == 0 {
		return components.MutedText.Render("No parameters")
	}
	rows := make([]string, 0, len(function.Parameters))
	for _, parameter := range function.Parameters {
		rows = append(rows, components.KeyValue(parameter.Name+" ("+parameter.Type+")", parameter.Value, 72))
	}
	return strings.Join(rows, "\n")
}
