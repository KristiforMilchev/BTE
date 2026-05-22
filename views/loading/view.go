package loading

import "bos/utils"

func (m *Model) View() string {
	return utils.RenderCentered(
		m.width,
		m.height,
		"Connecting to Ledger",
		"Requirements:\n- Ledger plugged in\n- Device unlocked\n- Ethereum app open\n- Ledger Live closed\n\nPress q to quit.",
	)
}
