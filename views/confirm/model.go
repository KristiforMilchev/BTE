package confirm

import (
	"bos/interfaces"
	"bos/views"
)

type Model struct {
	wallet interfaces.IWallet
	draft  views.TxDraft
	width  int
	height int
}

func New(wallet interfaces.IWallet, draft views.TxDraft) *Model {
	return &Model{wallet: wallet, draft: draft}
}
