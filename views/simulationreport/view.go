package simulationreport

import (
	"fmt"
	"math/big"
	"strings"

	"bos/components"
	"bos/enums"
	"bos/types"
	"bos/utils"
	"bos/views"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	report := m.payload.Report
	width := components.Clamp(m.width-8, 82, 124)
	contentWidth := components.Max(1, width-components.PanelStyle.GetHorizontalFrameSize())

	sections := []string{
		summary(report, contentWidth),
		balanceSection(report.BalanceChanges, contentWidth),
		detailGrid(report, contentWidth),
	}
	if warnings := warningRows(report.Warnings); len(warnings) > 0 {
		sections = append(sections, section("Warnings", warnings, contentWidth, false))
	}
	sections = append(sections, footer(m.running))

	body := components.Panel(width, strings.Join(nonEmpty(sections), "\n\n"))
	return views.RenderApp(m.width, m.height, enums.FocusSend, "Simulation report", func(width, height int) string {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, body)
	})
}

func summary(report types.SimulationReport, width int) string {
	title := report.Title
	if title == "" {
		title = "Transaction Simulation"
	}

	status := report.Status
	if status == "" {
		status = "pending"
	}

	meta := lipgloss.JoinHorizontal(
		lipgloss.Top,
		pill("status", statusText(status, report.RiskLevel), halfPillWidth(width, 0)),
		"  ",
		pill("risk", riskText(report.RiskLevel), halfPillWidth(width, 1)),
	)

	return strings.Join(nonEmpty([]string{
		components.SectionTitle.Render(title),
		meta,
		components.MutedText.Render(components.Truncate(report.Summary, width)),
	}), "\n")
}

func detailGrid(report types.SimulationReport, width int) string {
	gap := 3
	leftWidth := components.Max(24, (width-gap)/2)
	rightWidth := components.Max(24, width-gap-leftWidth)
	bytecode := contractBytecodeChecks(report.BytecodeChecks)

	rows := []string{}

	approvals := section("Approvals", approvalRows(report.TokenApprovals, leftWidth), leftWidth, false)
	actions := section("Performed Actions", callRows(report.Calls, rightWidth), rightWidth, true)
	if approvals != "" || actions != "" {
		rows = append(rows, sectionPair(approvals, actions, leftWidth, rightWidth, gap))
	}

	checkedBytecode := section("Checked Bytecode", bytecodeRows(bytecode, leftWidth), leftWidth, false)
	events := section("Events", eventRows(report.Events, rightWidth), rightWidth, false)
	if checkedBytecode != "" || events != "" {
		rows = append(rows, sectionPair(checkedBytecode, events, leftWidth, rightWidth, gap))
	}

	return strings.Join(nonEmpty(rows), "\n\n")
}

func balanceSection(changes []types.BalanceChange, width int) string {
	return section("Balance Changes", balanceRows(changes, width), width, false)
}

func section(title string, rows []string, width int, emptyAsNone bool) string {
	if len(rows) == 0 {
		if !emptyAsNone {
			return ""
		}
		rows = []string{components.MutedText.Render("none")}
	}
	return components.SectionTitle.Render(title) + "\n" +
		components.Separator(width) + "\n" +
		strings.Join(rows, "\n")
}

func sectionPair(left string, right string, leftWidth int, rightWidth int, gap int) string {
	if strings.TrimSpace(left) == "" {
		return fitBlock(right, leftWidth+gap+rightWidth)
	}
	if strings.TrimSpace(right) == "" {
		return fitBlock(left, leftWidth+gap+rightWidth)
	}

	left = fitBlock(left, leftWidth)
	right = fitBlock(right, rightWidth)
	height := components.Max(lipgloss.Height(left), lipgloss.Height(right))

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(leftWidth).Height(height).Render(left),
		strings.Repeat(" ", gap),
		lipgloss.NewStyle().Width(rightWidth).Height(height).Render(right),
	)
}

func fitBlock(value string, width int) string {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		lines[i] = components.Truncate(line, width)
	}
	return strings.Join(lines, "\n")
}

func halfPillWidth(width int, index int) int {
	available := components.Max(1, width-2)
	base := available / 2
	if index == 0 && available%2 != 0 {
		return base + 1
	}
	return base
}

func pill(label string, value string, width int) string {
	innerWidth := components.Max(1, width-4)
	labelWidth := components.Min(10, components.Max(6, innerWidth/3))
	valueWidth := components.Max(1, innerWidth-labelWidth-1)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		components.Label.Width(labelWidth).Render(label),
		" ",
		lipgloss.NewStyle().Width(valueWidth).Align(lipgloss.Center).Render(value),
	)

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(components.BorderOverlayer).
		Padding(0, 1).
		Width(innerWidth).
		MaxWidth(innerWidth)

	return style.Render(content)
}

func statusText(status string, risk string) string {
	if strings.EqualFold(risk, "high") || strings.EqualFold(risk, "critical") {
		return components.ErrorText.Render(status)
	}
	if strings.EqualFold(risk, "medium") || strings.EqualFold(risk, "pending") {
		return lipgloss.NewStyle().Foreground(components.Warning).Bold(true).Render(status)
	}
	return components.SuccessText.Render(status)
}

func riskText(risk string) string {
	if risk == "" {
		risk = "low"
	}
	if strings.EqualFold(risk, "high") || strings.EqualFold(risk, "critical") {
		return components.ErrorText.Render(risk)
	}
	if strings.EqualFold(risk, "medium") || strings.EqualFold(risk, "pending") {
		return lipgloss.NewStyle().Foreground(components.Warning).Bold(true).Render(risk)
	}
	return components.SuccessText.Render(risk)
}

func balanceRows(changes []types.BalanceChange, width int) []string {
	rows := make([]string, 0, len(changes)*2)
	for _, change := range changes {
		line := fmt.Sprintf("%s  %s -> %s", change.Asset, change.Before, change.After)
		rows = append(rows, components.Value.Render(components.Truncate(line, width)))
		if change.Delta != "" {
			rows = append(rows, components.MutedText.Render(components.Truncate("delta "+change.Delta, width)))
		}
	}
	return rows
}

func approvalRows(approvals []types.TokenApproval, width int) []string {
	rows := make([]string, 0, len(approvals)*2)
	for _, approval := range approvals {
		rows = append(rows, components.Value.Render(components.Truncate(approval.Token+"  "+approval.Amount, width)))
		rows = append(rows, components.MutedText.Render(components.Truncate("spender "+components.ShortAddress(approval.Spender)+"  risk "+approval.Risk, width)))
	}
	return rows
}

func bytecodeRows(checks []types.BytecodeCheck, width int) []string {
	rows := make([]string, 0, len(checks)*4)
	for _, check := range checks {
		kind := "wallet"
		if check.IsContract {
			kind = "contract"
		}
		rows = append(rows, components.Value.Render(components.ShortAddress(check.Address))+"  "+components.MutedText.Render(kind))
		rows = append(rows, components.MutedText.Render("hex "+components.Truncate(check.RuntimeHex, width-4)))
		rows = append(rows, components.MutedText.Render("bin "+components.Truncate(check.RuntimeBinary, width-4)))
		if check.Note != "" {
			rows = append(rows, components.MutedText.Render(components.Truncate(check.Note, width)))
		}
	}
	return rows
}

func contractBytecodeChecks(checks []types.BytecodeCheck) []types.BytecodeCheck {
	filtered := make([]types.BytecodeCheck, 0, len(checks))
	for _, check := range checks {
		if check.IsContract {
			filtered = append(filtered, check)
		}
	}
	return filtered
}

func callRows(calls []types.ContractCall, width int) []string {
	rows := make([]string, 0, len(calls)*2)
	for _, call := range calls {
		indent := strings.Repeat("  ", call.Depth)
		line := fmt.Sprintf("%s%s -> %s", indent, components.ShortAddress(call.From), components.ShortAddress(call.To))
		rows = append(rows, components.Value.Render(components.Truncate(line, width)))
		details := fmt.Sprintf("%s%s  value %s", indent, call.Function, call.Value)
		rows = append(rows, components.MutedText.Render(components.Truncate(details, width)))
	}
	return rows
}

func eventRows(events []types.EventLog, width int) []string {
	rows := make([]string, 0, len(events)*2)
	for _, event := range events {
		rows = append(rows, components.Value.Render(components.ShortAddress(event.Contract)+"  "+event.Name))
		if event.Details != "" {
			rows = append(rows, components.MutedText.Render(components.Truncate(event.Details, width)))
		}
	}
	return rows
}

func warningRows(warnings []string) []string {
	rows := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		normalized := strings.ToLower(strings.TrimSpace(warning))
		if normalized == "" ||
			normalized == "no authentication or rate limit checks are applied" ||
			normalized == "signed transaction was executed only on the ganache fork and was not broadcast to the upstream network" {
			continue
		}
		rows = append(rows, lipgloss.NewStyle().Foreground(components.Warning).Render(warning))
	}
	return rows
}

func footer(running bool) string {
	if running {
		return components.Button("waiting for Ledger signature", true) + " " + components.Button("esc  Back", false)
	}
	return components.Button("enter  Back to dashboard", true) + " " + components.Button("esc  Back", false)
}

func pendingBalanceAfter(draft types.TxDraft) string {
	before, err := utils.ParseEtherToWei(draft.Asset.Balance)
	if err != nil {
		return "pending"
	}
	amount, err := utils.ParseEtherToWei(draft.Amount)
	if err != nil {
		return "pending"
	}

	after := new(big.Int).Sub(before, amount)
	if after.Sign() < 0 {
		return "below zero"
	}
	return utils.WeiToEther(after)
}

func nonEmpty(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, value)
		}
	}
	return out
}
