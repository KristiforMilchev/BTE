package types

type SimulationReportPayload struct {
	Report SimulationReport `json:"report"`
	Draft  *TxDraft         `json:"draft,omitempty"`
	Return string           `json:"return,omitempty"`
}

type SimulationSession struct {
	ID              string              `json:"id"`
	RPC             string              `json:"rpc"`
	ChainID         string              `json:"chainId,omitempty"`
	Transaction     LedgerTransaction   `json:"transaction,omitempty"`
	Transactions    []LedgerTransaction `json:"transactions,omitempty"`
	Network         Network             `json:"network"`
	Address         string              `json:"address"`
	Caller          string              `json:"caller"`
	Amount          string              `json:"amount,omitempty"`
	Asset           string              `json:"asset,omitempty"`
	BalanceBefore   string              `json:"balanceBefore,omitempty"`
	BalanceAfter    string              `json:"balanceAfter,omitempty"`
	AddressContract bool                `json:"addressContract"`
	RuntimeHex      string              `json:"runtimeHex,omitempty"`
	RuntimeBinary   string              `json:"runtimeBinary,omitempty"`
}

type LedgerTransaction struct {
	Type      string `json:"type"`
	ChainID   string `json:"chainId"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Data      string `json:"data"`
	Nonce     string `json:"nonce"`
	Gas       string `json:"gas"`
	GasPrice  string `json:"gasPrice,omitempty"`
	Function  string `json:"function,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type SimulationFunctionCall struct {
	Function     string   `json:"function"`
	Signature    string   `json:"signature"`
	PassedData   string   `json:"passedData"`
	Status       string   `json:"status"`
	Consequences []string `json:"consequences"`
}

type SignedSimulationTransaction struct {
	Function          string `json:"function,omitempty"`
	Signature         string `json:"signature,omitempty"`
	SignedTransaction string `json:"signedTransaction"`
}

type SimulationReport struct {
	Title          string                   `json:"title"`
	Status         string                   `json:"status"`
	RiskLevel      string                   `json:"riskLevel"`
	GasEstimate    string                   `json:"gasEstimate"`
	Summary        string                   `json:"summary"`
	BalanceChanges []BalanceChange          `json:"balanceChanges"`
	TokenApprovals []TokenApproval          `json:"tokenApprovals"`
	BytecodeChecks []BytecodeCheck          `json:"bytecodeChecks"`
	Calls          []ContractCall           `json:"calls"`
	FunctionCalls  []SimulationFunctionCall `json:"functionCalls,omitempty"`
	Events         []EventLog               `json:"events"`
	Warnings       []string                 `json:"warnings"`
}

type BalanceChange struct {
	Asset  string `json:"asset"`
	Before string `json:"before"`
	After  string `json:"after"`
	Delta  string `json:"delta"`
}

type TokenApproval struct {
	Token   string `json:"token"`
	Spender string `json:"spender"`
	Amount  string `json:"amount"`
	Risk    string `json:"risk"`
}

type BytecodeCheck struct {
	Address       string `json:"address"`
	IsContract    bool   `json:"isContract"`
	RuntimeHex    string `json:"runtimeHex"`
	RuntimeBinary string `json:"runtimeBinary"`
	Note          string `json:"note"`
}

type ContractCall struct {
	Depth    int    `json:"depth"`
	From     string `json:"from"`
	To       string `json:"to"`
	Function string `json:"function"`
	Value    string `json:"value"`
}

type EventLog struct {
	Contract string `json:"contract"`
	Name     string `json:"name"`
	Details  string `json:"details"`
}

type ContractFunction struct {
	Name       string              `json:"name"`
	Signature  string              `json:"signature"`
	Mutability string              `json:"mutability"`
	Parameters []ContractParameter `json:"parameters"`
}

type ContractParameter struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ContractInteractionPayload struct {
	Address   string             `json:"address"`
	Functions []ContractFunction `json:"functions"`
}
