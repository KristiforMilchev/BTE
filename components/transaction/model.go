package transaction

import "bos/types"

type Model struct {
	transaction types.Transaction
	selected    bool
}

func New(tx types.Transaction) Model {
	return Model{transaction: tx}
}

func (m *Model) SetSelected(selected bool) {
	m.selected = selected
}
