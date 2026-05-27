package implementations

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"bos/constants"
	"bos/types"

	"github.com/ethereum/go-ethereum/common"
)

const (
	beginSimulationPath   = "/v1/simulation/begin"
	performSimulationPath = "/v1/simulation/perform"
)

type Simulator struct {
	baseURL string
	client  *http.Client
}

type beginSimulationRequest struct {
	Network   string                   `json:"network"`
	Address   string                   `json:"address"`
	Caller    string                   `json:"caller"`
	Functions []types.ContractFunction `json:"functions,omitempty"`
}

type beginSimulationResponse struct {
	SimulationID string                    `json:"simulationId"`
	ID           string                    `json:"id"`
	SessionID    string                    `json:"sessionId"`
	RPC          string                    `json:"rpc"`
	ClonedRPC    string                    `json:"clonedRpc"`
	ChainID      string                    `json:"chainId"`
	ChainIDAlt   string                    `json:"chainID"`
	Transaction  types.LedgerTransaction   `json:"transaction"`
	Transactions []types.LedgerTransaction `json:"transactions"`
	Contract     struct {
		HasCode bool   `json:"hasCode"`
		Address string `json:"address"`
	} `json:"contract"`
	AddressContract bool   `json:"addressContract"`
	IsContract      bool   `json:"isContract"`
	Bytecode        string `json:"bytecode"`
	RuntimeHex      string `json:"runtimeHex"`
	RuntimeBinary   string `json:"runtimeBinary"`
}

type executeSimulationRequest struct {
	SimulationID       string                              `json:"simulationId"`
	Session            types.SimulationSession             `json:"session"`
	SignedTx           string                              `json:"signedTx,omitempty"`
	SignedTransaction  string                              `json:"signedTransaction,omitempty"`
	RawTransaction     string                              `json:"rawTransaction,omitempty"`
	SignedTransactions []types.SignedSimulationTransaction `json:"signedTransactions,omitempty"`
}

type executeSimulationResponse struct {
	Report types.SimulationReport `json:"report"`
}

type performSimulationResponse struct {
	SimulationID         string               `json:"simulationId"`
	Network              string               `json:"network"`
	RawTransactionSHA256 string               `json:"rawTransactionSha256"`
	Balances             balanceSnapshot      `json:"balances"`
	ApprovalFindings     []approvalFinding    `json:"approvalFindings"`
	Contract             contractSnapshot     `json:"contract"`
	Execution            executionReport      `json:"execution"`
	FunctionCalls        []functionCallReport `json:"functionCalls,omitempty"`
	Warnings             []string             `json:"warnings,omitempty"`
	PerformedAt          time.Time            `json:"performedAt"`
}

type balanceSnapshot struct {
	CallerBefore  string `json:"callerBefore"`
	CallerAfter   string `json:"callerAfter,omitempty"`
	AddressBefore string `json:"addressBefore"`
	AddressAfter  string `json:"addressAfter,omitempty"`
}

type approvalFinding struct {
	Type        string `json:"type"`
	Selector    string `json:"selector"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type contractSnapshot struct {
	Address        string   `json:"address"`
	HasCode        bool     `json:"hasCode"`
	CodeHashSHA256 string   `json:"codeHashSha256,omitempty"`
	StateChanges   []string `json:"stateChanges,omitempty"`
}

type executionReport struct {
	Mode              string `json:"mode"`
	Broadcasted       bool   `json:"broadcasted"`
	TransactionHash   string `json:"transactionHash,omitempty"`
	Status            string `json:"status"`
	Details           string `json:"details"`
	ForkBackendNeeded bool   `json:"forkBackendNeeded"`
}

type functionCallReport struct {
	Function     string   `json:"function"`
	Signature    string   `json:"signature"`
	PassedData   string   `json:"passedData"`
	Status       string   `json:"status"`
	Consequences []string `json:"consequences"`
}

func (s *Simulator) BeginSimulation(network types.Network, address common.Address, caller common.Address) (*types.SimulationSession, error) {
	return s.beginSimulation(network, address, caller, nil)
}

func (s *Simulator) BeginContractSimulation(network types.Network, address common.Address, caller common.Address, functions []types.ContractFunction) (*types.SimulationSession, error) {
	return s.beginSimulation(network, address, caller, functions)
}

func (s *Simulator) beginSimulation(network types.Network, address common.Address, caller common.Address, functions []types.ContractFunction) (*types.SimulationSession, error) {
	response := beginSimulationResponse{}
	if err := s.postJSON(beginSimulationPath, beginSimulationRequest{
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

func (s *Simulator) ExecuteSignedTransaction(session types.SimulationSession, signedTx []byte) (*types.SimulationReport, error) {
	if len(signedTx) == 0 {
		return nil, fmt.Errorf("signed transaction is empty")
	}

	rawTx := "0x" + hex.EncodeToString(signedTx)
	responseBody, err := s.postJSONRaw(performSimulationPath, executeSimulationRequest{
		SimulationID:      session.ID,
		Session:           session,
		SignedTx:          rawTx,
		SignedTransaction: rawTx,
		RawTransaction:    rawTx,
	})
	if err != nil {
		return nil, err
	}

	var wrapped executeSimulationResponse
	_ = json.Unmarshal(responseBody, &wrapped)
	report := wrapped.Report
	if reportEmpty(report) {
		var apiResponse performSimulationResponse
		if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
			return nil, fmt.Errorf("decode simulator API response %s: %w", performSimulationPath, err)
		}
		if apiResponse.Execution.Status != "" || apiResponse.SimulationID != "" {
			report = reportFromPerformResponse(apiResponse, session)
		} else {
			var direct types.SimulationReport
			if err := json.Unmarshal(responseBody, &direct); err != nil {
				return nil, fmt.Errorf("decode simulator API response %s: %w", performSimulationPath, err)
			}
			report = direct
		}
	}
	applyReportFallbacks(&report, session)
	return &report, nil
}

func (s *Simulator) ExecuteSignedTransactions(session types.SimulationSession, signedTxs []types.SignedSimulationTransaction) (*types.SimulationReport, error) {
	if len(signedTxs) == 0 {
		return nil, fmt.Errorf("signed transactions are empty")
	}

	responseBody, err := s.postJSONRaw(performSimulationPath, executeSimulationRequest{
		SimulationID:       session.ID,
		Session:            session,
		SignedTransactions: signedTxs,
	})
	if err != nil {
		return nil, err
	}

	var apiResponse performSimulationResponse
	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return nil, fmt.Errorf("decode simulator API response %s: %w", performSimulationPath, err)
	}

	report := reportFromPerformResponse(apiResponse, session)
	applyReportFallbacks(&report, session)
	return &report, nil
}

func NewSimulator() *Simulator {
	return NewSimulatorWithBaseURL(constants.SimulatorBaseURL)
}

func NewSimulatorWithBaseURL(baseURL string) *Simulator {
	return NewSimulatorWithHTTPClient(baseURL, &http.Client{Timeout: 60 * time.Second})
}

func NewSimulatorWithHTTPClient(baseURL string, client *http.Client) *Simulator {
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}
	return &Simulator{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}

func (s *Simulator) postJSON(path string, payload any, target any) error {
	responseBody, err := s.postJSONRaw(path, payload)
	if err != nil {
		return err
	}
	if len(responseBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(responseBody, target); err != nil {
		return fmt.Errorf("decode simulator API response %s: %w", path, err)
	}
	return nil
}

func (s *Simulator) postJSONRaw(path string, payload any) ([]byte, error) {
	if s.client == nil {
		s.client = &http.Client{Timeout: 60 * time.Second}
	}
	if s.baseURL == "" {
		s.baseURL = constants.SimulatorBaseURL
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("simulator API %s: %w", path, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read simulator API response %s: %w", path, err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("simulator API %s returned %s: %s", path, resp.Status, strings.TrimSpace(string(responseBody)))
	}
	return responseBody, nil
}

func networkRPC(network types.Network) string {
	if network.Rpc != nil {
		return *network.Rpc
	}
	return ""
}

func applyReportFallbacks(report *types.SimulationReport, session types.SimulationSession) {
	if report.Title == "" {
		report.Title = "Transaction Simulation"
	}
	if report.Status == "" {
		report.Status = "Simulated"
	}
	if report.RiskLevel == "" {
		report.RiskLevel = "Low"
	}
	if report.Summary == "" {
		report.Summary = "Signed transaction was tested before broadcast."
	}
	if len(report.BalanceChanges) == 0 {
		report.BalanceChanges = []types.BalanceChange{
			{Asset: sessionAsset(session), Before: session.BalanceBefore, After: session.BalanceAfter, Delta: "-" + session.Amount},
		}
	}
	if len(report.BytecodeChecks) == 0 {
		report.BytecodeChecks = []types.BytecodeCheck{bytecodeCheck(session)}
	}
	if len(report.Calls) == 0 {
		report.Calls = performedActions(session)
	}
	report.Warnings = meaningfulWarnings(report.Warnings)
}

func reportFromPerformResponse(response performSimulationResponse, session types.SimulationSession) types.SimulationReport {
	report := types.SimulationReport{
		Title:          simulationTitle(session),
		Status:         simulationStatus(response.Execution.Status),
		RiskLevel:      riskLevel(response),
		GasEstimate:    session.Transaction.Gas,
		Summary:        simulationSummary(response, session),
		BalanceChanges: balanceChanges(response, session),
		TokenApprovals: tokenApprovals(response.ApprovalFindings),
		BytecodeChecks: []types.BytecodeCheck{
			{
				Address:    firstNonEmpty(response.Contract.Address, session.Address),
				IsContract: response.Contract.HasCode,
				Note:       codeHashNote(response.Contract.CodeHashSHA256),
			},
		},
		Calls: []types.ContractCall{
			{
				Depth:    0,
				From:     firstNonEmpty(session.Transaction.From, session.Caller),
				To:       firstNonEmpty(session.Transaction.To, session.Address),
				Function: simulationAction(response, session),
				Value:    simulationValue(session),
			},
		},
		FunctionCalls: functionCalls(response.FunctionCalls),
		Events:        eventsFromPerformResponse(response),
		Warnings:      response.Warnings,
	}

	if response.Execution.TransactionHash != "" {
		report.Events = append(report.Events, types.EventLog{
			Contract: firstNonEmpty(response.Contract.Address, session.Address),
			Name:     "Transaction Hash",
			Details:  response.Execution.TransactionHash,
		})
	}
	return report
}

func simulationTitle(session types.SimulationSession) string {
	if session.AddressContract {
		return "Contract Simulation"
	}
	return "Transaction Simulation"
}

func simulationSummary(response performSimulationResponse, session types.SimulationSession) string {
	if response.Execution.Details != "" {
		return response.Execution.Details
	}
	if session.AddressContract {
		return "Function simulations completed before broadcast."
	}
	return "Native transfer simulation completed before broadcast."
}

func simulationAction(response performSimulationResponse, session types.SimulationSession) string {
	if session.AddressContract {
		return firstNonEmpty(response.Execution.Mode, "contract simulation")
	}
	return "regular wallet transfer"
}

func simulationValue(session types.SimulationSession) string {
	if session.Amount != "" {
		return session.Amount + " " + sessionAsset(session)
	}
	return firstNonEmpty(session.Transaction.Value, "0x0")
}

func functionCalls(calls []functionCallReport) []types.SimulationFunctionCall {
	out := make([]types.SimulationFunctionCall, 0, len(calls))
	for _, call := range calls {
		out = append(out, types.SimulationFunctionCall{
			Function:     call.Function,
			Signature:    call.Signature,
			PassedData:   call.PassedData,
			Status:       simulationStatus(call.Status),
			Consequences: cleanConsequences(call.Consequences),
		})
	}
	return out
}

func cleanConsequences(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" ||
			strings.Contains(normalized, "debug_tracetransaction") ||
			strings.Contains(normalized, ".sha256=") {
			continue
		}
		out = append(out, value)
	}
	return out
}

func balanceChanges(response performSimulationResponse, session types.SimulationSession) []types.BalanceChange {
	changes := make([]types.BalanceChange, 0, 2)
	if response.Balances.CallerBefore != "" || response.Balances.CallerAfter != "" {
		changes = append(changes, types.BalanceChange{
			Asset:  networkSymbol(session.Network),
			Before: response.Balances.CallerBefore,
			After:  response.Balances.CallerAfter,
			Delta:  balanceDelta(response.Balances.CallerBefore, response.Balances.CallerAfter),
		})
	}
	if response.Balances.AddressBefore != "" || response.Balances.AddressAfter != "" {
		changes = append(changes, types.BalanceChange{
			Asset:  recipientBalanceAsset(session),
			Before: response.Balances.AddressBefore,
			After:  response.Balances.AddressAfter,
			Delta:  balanceDelta(response.Balances.AddressBefore, response.Balances.AddressAfter),
		})
	}
	return changes
}

func recipientBalanceAsset(session types.SimulationSession) string {
	if session.AddressContract {
		return "Contract native balance"
	}
	return "Recipient native balance"
}

func tokenApprovals(findings []approvalFinding) []types.TokenApproval {
	approvals := make([]types.TokenApproval, 0, len(findings))
	for _, finding := range findings {
		approvals = append(approvals, types.TokenApproval{
			Token:   firstNonEmpty(finding.Type, "approval"),
			Spender: firstNonEmpty(finding.Selector, "unknown"),
			Amount:  finding.Description,
			Risk:    firstNonEmpty(finding.Severity, "medium"),
		})
	}
	return approvals
}

func eventsFromPerformResponse(response performSimulationResponse) []types.EventLog {
	events := make([]types.EventLog, 0, len(response.Contract.StateChanges)+1)
	for _, change := range response.Contract.StateChanges {
		events = append(events, types.EventLog{
			Contract: response.Contract.Address,
			Name:     "State Change",
			Details:  change,
		})
	}
	return events
}

func riskLevel(response performSimulationResponse) string {
	level := "Low"
	for _, finding := range response.ApprovalFindings {
		switch strings.ToLower(finding.Severity) {
		case "critical":
			return "Critical"
		case "high":
			level = "High"
		case "medium":
			if level == "Low" {
				level = "Medium"
			}
		}
	}
	if response.Execution.Status != "" && response.Execution.Status != "0x1" && !strings.EqualFold(response.Execution.Status, "success") {
		if level == "Low" {
			level = "Medium"
		}
	}
	return level
}

func simulationStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "0x1", "1", "success", "passed":
		return "passed"
	case "0x0", "0", "failed", "reverted":
		return "reverted"
	case "":
		return "simulated"
	default:
		return status
	}
}

func codeHashNote(codeHash string) string {
	if codeHash == "" {
		return "Contract bytecode proof captured."
	}
	return "Bytecode proof sha256=" + codeHash
}

func balanceDelta(before string, after string) string {
	if before == "" || after == "" {
		return ""
	}
	if before == after {
		return "0x0"
	}
	return before + " -> " + after
}

func reportEmpty(report types.SimulationReport) bool {
	return report.Title == "" &&
		report.Status == "" &&
		report.RiskLevel == "" &&
		report.GasEstimate == "" &&
		report.Summary == "" &&
		len(report.BalanceChanges) == 0 &&
		len(report.TokenApprovals) == 0 &&
		len(report.BytecodeChecks) == 0 &&
		len(report.Calls) == 0 &&
		len(report.Events) == 0 &&
		len(report.Warnings) == 0
}

func networkSymbol(network types.Network) string {
	if network.Symbol == nil || *network.Symbol == "" {
		return "ETH"
	}
	return *network.Symbol
}

func sessionAsset(session types.SimulationSession) string {
	if session.Asset != "" {
		return session.Asset
	}
	return networkSymbol(session.Network)
}

func performedActions(session types.SimulationSession) []types.ContractCall {
	if !session.AddressContract {
		return []types.ContractCall{
			{
				Depth:    0,
				From:     session.Caller,
				To:       session.Address,
				Function: "regular wallet transfer",
				Value:    session.Amount + " " + sessionAsset(session),
			},
		}
	}

	return []types.ContractCall{
		{Depth: 0, From: session.Caller, To: session.Address, Function: "contract call", Value: session.Amount + " " + sessionAsset(session)},
	}
}

func regularWalletNote(session types.SimulationSession) string {
	if session.AddressContract {
		return "Recipient has contract bytecode in the simulator report."
	}
	return "Recipient is not a contract in the simulator report; simulated as a regular wallet transfer."
}

func meaningfulWarnings(warnings []string) []string {
	meaningful := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		normalized := strings.ToLower(strings.TrimSpace(warning))
		if normalized == "" {
			continue
		}
		if normalized == "no authentication or rate limit checks are applied" {
			continue
		}
		if normalized == "signed transaction was executed only on the ganache fork and was not broadcast to the upstream network" {
			continue
		}
		meaningful = append(meaningful, warning)
	}
	return meaningful
}

func bytecodeCheck(session types.SimulationSession) types.BytecodeCheck {
	if !session.AddressContract {
		return types.BytecodeCheck{
			Address:       session.Address,
			IsContract:    false,
			RuntimeHex:    "0x",
			RuntimeBinary: "0b0",
			Note:          "No runtime bytecode returned for this address.",
		}
	}

	runtimeHex := firstNonEmpty(session.RuntimeHex, "0x")
	runtimeBinary := firstNonEmpty(session.RuntimeBinary, "0b1")
	return types.BytecodeCheck{
		Address:       session.Address,
		IsContract:    true,
		RuntimeHex:    runtimeHex,
		RuntimeBinary: runtimeBinary,
		Note:          "Contract bytecode proof captured.",
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
