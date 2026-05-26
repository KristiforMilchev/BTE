package importcontract

import (
	"bos/interfaces"
	"bos/types"

	"github.com/charmbracelet/bubbles/textinput"
)

type Method struct {
	Name      string
	Signature string
}

type Model struct {
	input               textinput.Model
	reader              interfaces.IContractInteractionReader
	address             string
	loaded              bool
	status              string
	callable            []Method
	readable            []Method
	interactions        []types.ContractInteraction
	loadingInteractions bool
	simulating          bool
	width               int
	height              int
}

func New(readers ...interfaces.IContractInteractionReader) *Model {
	input := textinput.New()
	input.Placeholder = "Contract address"
	input.CharLimit = 128
	input.Width = 58
	input.Focus()

	var reader interfaces.IContractInteractionReader
	if len(readers) > 0 {
		reader = readers[0]
	}

	return &Model{
		input:  input,
		reader: reader,
		status: "No contract address selected",
	}
}

func fakeCallableMethods() []Method {
	return []Method{
		{Name: "transfer", Signature: "transfer(address,uint256)"},
		{Name: "approve", Signature: "approve(address,uint256)"},
		{Name: "transferFrom", Signature: "transferFrom(address,address,uint256)"},
	}
}

func fakeReadableMethods() []Method {
	return []Method{
		{Name: "name", Signature: "name()"},
		{Name: "symbol", Signature: "symbol()"},
		{Name: "decimals", Signature: "decimals()"},
		{Name: "balanceOf", Signature: "balanceOf(address)"},
		{Name: "allowance", Signature: "allowance(address,address)"},
		{Name: "totalSupply", Signature: "totalSupply()"},
	}
}
