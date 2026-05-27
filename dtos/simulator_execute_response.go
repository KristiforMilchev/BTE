package dtos

import "bos/types"

type SimulatorExecuteResponse struct {
	Report types.SimulationReport `json:"report"`
}
