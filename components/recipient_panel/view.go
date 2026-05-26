package recipientPanel

import (
	"bos/components"
	"bos/types"
	"strings"
)

func Render(recipient types.Contact, width int) string {
	return strings.Join([]string{
		components.SectionTitle.Render("Recipient"),
		components.Value.Render(components.Truncate(recipient.Name, width)),
		components.MutedText.Render(components.ShortAddress(recipient.Address)),
	}, "\n")
}
