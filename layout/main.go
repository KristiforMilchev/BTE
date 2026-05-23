package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Constraints struct {
	Width  int
	Height int
}

type Context struct {
	Constraints Constraints
}

type Widget interface {
	Render(ctx Context) string
}

type WidgetFunc func(ctx Context) string

func (f WidgetFunc) Render(ctx Context) string {
	return f(ctx)
}

type Child struct {
	Widget Widget
	Fixed  int
	Weight int
}

func View(fn func(width int, height int) string) Widget {
	return WidgetFunc(func(ctx Context) string {
		return fn(ctx.Constraints.Width, ctx.Constraints.Height)
	})
}

func Text(value string) Widget {
	return WidgetFunc(func(ctx Context) string {
		return lipgloss.NewStyle().
			Width(ctx.Constraints.Width).
			Height(ctx.Constraints.Height).
			Render(value)
	})
}

func Fixed(size int, widget Widget) Child {
	return Child{Fixed: size, Widget: widget}
}

func Expanded(weight int, widget Widget) Child {
	if weight <= 0 {
		weight = 1
	}

	return Child{Weight: weight, Widget: widget}
}

func Gap(size int) Child {
	return Fixed(size, WidgetFunc(func(ctx Context) string {
		height := max(1, ctx.Constraints.Height)
		lines := make([]string, height)

		for i := range lines {
			lines[i] = strings.Repeat(" ", max(0, size))
		}

		return strings.Join(lines, "\n")
	}))
}

func Row(children ...Child) Widget {
	return WidgetFunc(func(ctx Context) string {
		width := max(0, ctx.Constraints.Width)
		height := max(0, ctx.Constraints.Height)

		fixed := 0
		totalWeight := 0

		for _, child := range children {
			if child.Fixed > 0 {
				fixed += child.Fixed
				continue
			}

			totalWeight += max(1, child.Weight)
		}

		remaining := max(0, width-fixed)
		used := 0
		parts := make([]string, 0, len(children))

		for i, child := range children {
			childWidth := child.Fixed

			if child.Fixed <= 0 {
				weight := max(1, child.Weight)
				childWidth = remaining * weight / max(1, totalWeight)

				// Give rounding remainder to the last flexible child so rows fill exactly.
				if isLastFlexible(children, i) {
					childWidth = max(0, remaining-used)
				}
			}

			used += childWidth

			if child.Widget == nil {
				parts = append(parts, "")
				continue
			}

			parts = append(parts, child.Widget.Render(Context{
				Constraints: Constraints{Width: childWidth, Height: height},
			}))
		}

		return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	})
}

func Column(children ...Child) Widget {
	return WidgetFunc(func(ctx Context) string {
		width := max(0, ctx.Constraints.Width)
		height := max(0, ctx.Constraints.Height)

		fixed := 0
		totalWeight := 0

		for _, child := range children {
			if child.Fixed > 0 {
				fixed += child.Fixed
				continue
			}

			totalWeight += max(1, child.Weight)
		}

		remaining := max(0, height-fixed)
		used := 0
		parts := make([]string, 0, len(children))

		for i, child := range children {
			childHeight := child.Fixed

			if child.Fixed <= 0 {
				weight := max(1, child.Weight)
				childHeight = remaining * weight / max(1, totalWeight)

				if isLastFlexible(children, i) {
					childHeight = max(0, remaining-used)
				}
			}

			used += childHeight

			if child.Widget == nil {
				parts = append(parts, "")
				continue
			}

			parts = append(parts, child.Widget.Render(Context{
				Constraints: Constraints{Width: width, Height: childHeight},
			}))
		}

		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	})
}

func Padding(x int, y int, child Widget) Widget {
	return WidgetFunc(func(ctx Context) string {
		width := max(0, ctx.Constraints.Width)
		height := max(0, ctx.Constraints.Height)

		innerWidth := max(0, width-(x*2))
		innerHeight := max(0, height-(y*2))

		body := ""
		if child != nil {
			body = child.Render(Context{
				Constraints: Constraints{Width: innerWidth, Height: innerHeight},
			})
		}

		return lipgloss.NewStyle().
			Width(width).
			Height(height).
			Padding(y, x).
			Render(body)
	})
}

func Center(child Widget) Widget {
	return WidgetFunc(func(ctx Context) string {
		if child == nil {
			return ""
		}
		rendered := child.Render(ctx)
		return lipgloss.Place(ctx.Constraints.Width, ctx.Constraints.Height, lipgloss.Center, lipgloss.Center, rendered)
	})
}

func CenterH(child Widget) Widget {
	return WidgetFunc(func(ctx Context) string {
		if child == nil {
			return ""
		}
		rendered := child.Render(ctx)
		return lipgloss.PlaceHorizontal(ctx.Constraints.Width, lipgloss.Center, rendered)
	})
}

func Constrain(width int, height int, child Widget) Widget {
	return WidgetFunc(func(ctx Context) string {
		if child == nil {
			return ""
		}

		return child.Render(Context{
			Constraints: Constraints{
				Width:  choose(width, ctx.Constraints.Width),
				Height: choose(height, ctx.Constraints.Height),
			},
		})
	})
}

func Render(width int, height int, widget Widget) string {
	if widget == nil {
		return ""
	}

	return widget.Render(Context{
		Constraints: Constraints{Width: width, Height: height},
	})
}

func isLastFlexible(children []Child, index int) bool {
	for i := index + 1; i < len(children); i++ {
		if children[i].Fixed <= 0 {
			return false
		}
	}

	return true
}

func choose(value int, fallback int) int {
	if value > 0 {
		return value
	}

	return fallback
}

