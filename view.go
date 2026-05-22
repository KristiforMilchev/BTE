package main

import (
	"bos/constants"
	"bos/views"
)

func (m Model) View() string {
	return views.Render(views.State{
		Screen:           m.screen,
		Focus:            m.focus,
		Width:            m.width,
		Height:           m.height,
		Address:          m.address,
		Balance:          m.balance,
		ChainID:          m.chainID,
		Err:              m.err,
		TxHash:           m.txHash,
		AmountInput:      m.amountInput.View(),
		AmountValue:      m.amountInput.Value(),
		SelectedToken:    m.selectedToken,
		SelectedContact:  m.selectedContact,
		Tokens:           m.tokens,
		Contacts:         m.contacts,
		SimulationStatus: m.simulationStatus,
		RiskLevel:        m.riskLevel,
		EstimatedFee:     m.estimatedFee,
		StatusMessage:    m.statusMessage,
		RpcURL:           constants.RpcURL,
	})
}
