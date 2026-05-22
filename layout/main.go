package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Node interface {
	Render(width int, height int) string
}

type RenderFunc func(width int, height int) string

func (f RenderFunc) Render(width int, height int) string {
	return f(width, height)
}

type Child struct {
	Node     Node
	Fixed    int
	Expanded bool
}

func Fixed(size int, node Node) Child {
	return Child{Node: node, Fixed: size}
}

func Expanded(node Node) Child {
	return Child{Node: node, Expanded: true}
}

func View(fn func(width int, height int) string) Node {
	return RenderFunc(fn)
}

func Text(value string) Node {
	return RenderFunc(func(width int, height int) string {
		return lipgloss.NewStyle().
			Width(width).
			Height(height).
			Render(value)
	})
}

func Row(children ...Child) Node {
	return RenderFunc(func(width int, height int) string {
		fixed := 0
		expandedCount := 0

		for _, child := range children {
			if child.Expanded {
				expandedCount++
			} else {
				fixed += child.Fixed
			}
		}

		remaining := width - fixed
		if remaining < 0 {
			remaining = 0
		}

		expandedWidth := 0
		if expandedCount > 0 {
			expandedWidth = remaining / expandedCount
		}

		parts := make([]string, 0, len(children))

		for _, child := range children {
			childWidth := child.Fixed
			if child.Expanded {
				childWidth = expandedWidth
			}

			parts = append(parts, child.Node.Render(childWidth, height))
		}

		return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	})
}

func Col(children ...Child) Node {
	return RenderFunc(func(width int, height int) string {
		fixed := 0
		expandedCount := 0

		for _, child := range children {
			if child.Expanded {
				expandedCount++
			} else {
				fixed += child.Fixed
			}
		}

		remaining := height - fixed
		if remaining < 0 {
			remaining = 0
		}

		expandedHeight := 0
		if expandedCount > 0 {
			expandedHeight = remaining / expandedCount
		}

		parts := make([]string, 0, len(children))

		for _, child := range children {
			childHeight := child.Fixed
			if child.Expanded {
				childHeight = expandedHeight
			}

			parts = append(parts, child.Node.Render(width, childHeight))
		}

		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	})
}

func Gap(size int) Node {
	return RenderFunc(func(width int, height int) string {
		lines := make([]string, max(1, height))
		for i := range lines {
			lines[i] = strings.Repeat(" ", size)
		}
		return strings.Join(lines, "\n")
	})
}

func Panel(node Node) Node {
	return RenderFunc(func(width int, height int) string {
		innerWidth := max(0, width-2)
		innerHeight := max(0, height-2)

		return lipgloss.NewStyle().
			Width(width).
			Height(height).
			Border(lipgloss.NormalBorder()).
			Render(node.Render(innerWidth, innerHeight))
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
