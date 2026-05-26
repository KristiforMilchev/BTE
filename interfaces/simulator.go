package interfaces

import (
	"bos/types"

	"github.com/ethereum/go-ethereum/common"
)

type ISimulator interface {
	BeginSimulation(network types.Network, address common.Address, caller common.Address) (*types.SimulationSession, error)
	ExecuteSignedTransaction(session types.SimulationSession, signedTx []byte) (*types.SimulationReport, error)
}
