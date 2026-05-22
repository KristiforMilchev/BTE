package footer

import (
	"bos/components"
	"bos/utils"
	"strings"
)

func RenderFooter(width int, statusMessage string, rpcURL string) string {
	width = utils.SafeWidth(width)

	status := "Ledger Connected • RPC Online"
	if rpcURL != "" {
		status += " • " + rpcURL
	}
	if strings.TrimSpace(statusMessage) != "" {
		status += " • " + statusMessage
	}

	return components.Separator(width) + "\n" + components.MutedText.Render(components.Truncate(status, width))
}
