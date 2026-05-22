package panel

import (
	"strings"

	"bos/components"
	"bos/layout"

	"github.com/charmbracelet/lipgloss"
)

type Props struct {
	Title string
}

func New(title string, child layout.Widget) layout.Widget {
	return View(Props{Title: title}, child)
}

func View(props Props, child layout.Widget) layout.Widget {
	return layout.WidgetFunc(func(ctx layout.Context) string {
		width := components.Max(12, ctx.Constraints.Width)
		height := components.Max(1, ctx.Constraints.Height)

		frameWidth := components.PanelStyle.GetHorizontalFrameSize()
		frameHeight := components.PanelStyle.GetVerticalFrameSize()

		innerWidth := components.Max(1, width-frameWidth-4)
		innerHeight := components.Max(1, height-frameHeight-3)

		body := ""
		if child != nil {
			body = child.Render(layout.Context{
				Constraints: layout.Constraints{Width: innerWidth, Height: innerHeight},
			})
		}

		content := strings.Join([]string{
			components.SectionTitle.
				Width(innerWidth).
				AlignHorizontal(lipgloss.Center).
				Render(props.Title),
			"",
			body,
		}, "\n")

		return components.PanelSized(width, height, content)
	})
}
