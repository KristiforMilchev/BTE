package confirm

import (
	"bos/interfaces"
	"bos/types"
)

type Model struct {
	wallet interfaces.IWallet
	draft  types.TxDraft
	width  int
	height int
}

func New(wallet interfaces.IWallet, draft types.TxDraft) *Model {
	return &Model{wallet: wallet, draft: draft}
}
