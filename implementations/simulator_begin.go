package implementations

import (
	"fmt"

	"bos/dtos"
	"bos/types"

	"github.com/ethereum/go-ethereum/common"
)

func (s *Simulator) BeginSimulation(network types.Network, address common.Address, caller common.Address) (*types.SimulationSession, error) {
	return s.beginSimulation(network, address, caller, nil)
}

func (s *Simulator) BeginContractSimulation(network types.Network, address common.Address, caller common.Address, functions []types.ContractFunction) (*types.SimulationSession, error) {
	return s.beginSimulation(network, address, caller, functions)
}

func (s *Simulator) beginSimulation(network types.Network, address common.Address, caller common.Address, functions []types.ContractFunction) (*types.SimulationSession, error) {
	response := dtos.SimulatorBeginResponse{}
	if err := s.postJSON(beginSimulationPath, dtos.SimulatorBeginRequest{
		Network:   networkRPC(network),
		Address:   address.Hex(),
		Caller:    caller.Hex(),
		Functions: functions,
	}, &response); err != nil {
		return nil, err
	}

	rpc := firstNonEmpty(response.ClonedRPC, response.RPC)
	if rpc == "" {
		return nil, fmt.Errorf("simulator did not return an RPC")
	}

	id := response.SimulationID
	if id == "" {
		id = response.SessionID
	}
	if id == "" {
		id = response.ID
	}

	return &types.SimulationSession{
		ID:              id,
		RPC:             rpc,
		ChainID:         firstNonEmpty(response.Transaction.ChainID, response.ChainID, response.ChainIDAlt),
		Transaction:     response.Transaction,
		Transactions:    response.Transactions,
		Network:         network,
		Address:         address.Hex(),
		Caller:          caller.Hex(),
		AddressContract: response.AddressContract || response.IsContract || response.Contract.HasCode,
		RuntimeHex:      firstNonEmpty(response.RuntimeHex, response.Bytecode),
		RuntimeBinary:   response.RuntimeBinary,
	}, nil
}
