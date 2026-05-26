package networksetup

import (
	networkDialog "bos/components/network_dialog"
	"bos/interfaces"
	"bos/repositories"
)

type Model struct {
	network  interfaces.INetwork
	register repositories.RepositoryRegister
	dialog   *networkDialog.Model
	width    int
	height   int
}

func New(network interfaces.INetwork, register repositories.RepositoryRegister) *Model {
	dialog := networkDialog.New()
	dialog.Visible = true

	return &Model{
		network:  network,
		register: register,
		dialog:   dialog,
	}
}
