package contractinteraction

import "bos/types"

type Model struct {
	address   string
	functions []types.ContractFunction
	selected  int
	width     int
	height    int
}

func New(payload types.ContractInteractionPayload) *Model {
	if payload.Address == "" {
		payload.Address = "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC"
	}
	if len(payload.Functions) == 0 {
		payload.Functions = fakeFunctions()
	}

	return &Model{
		address:   payload.Address,
		functions: payload.Functions,
	}
}

func fakeFunctions() []types.ContractFunction {
	return []types.ContractFunction{
		{
			Name:       "transfer",
			Signature:  "transfer(address,uint256)",
			Mutability: "nonpayable",
			Parameters: []types.ContractParameter{
				{Name: "to", Type: "address", Value: "0x1111111111111111111111111111111111111111"},
				{Name: "amount", Type: "uint256", Value: "1000000000000000000"},
			},
		},
		{
			Name:       "approve",
			Signature:  "approve(address,uint256)",
			Mutability: "nonpayable",
			Parameters: []types.ContractParameter{
				{Name: "spender", Type: "address", Value: "0x2222222222222222222222222222222222222222"},
				{Name: "amount", Type: "uint256", Value: "unlimited"},
			},
		},
		{
			Name:       "balanceOf",
			Signature:  "balanceOf(address)",
			Mutability: "view",
			Parameters: []types.ContractParameter{
				{Name: "owner", Type: "address", Value: "0x3333333333333333333333333333333333333333"},
			},
		},
	}
}
