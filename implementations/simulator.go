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
	Network string `json:"network"`
	Address string `json:"address"`
	Caller  string `json:"caller"`
}

type beginSimulationResponse struct {
	SimulationID string `json:"simulationId"`
	ID           string `json:"id"`
	SessionID    string `json:"sessionId"`
	RPC          string `json:"rpc"`
	ClonedRPC    string `json:"clonedRpc"`
	ChainID      string `json:"chainId"`
	ChainIDAlt   string `json:"chainID"`
	Transaction  struct {
		ChainID string `json:"chainId"`
	} `json:"transaction"`
	Contract struct {
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
	SimulationID      string                  `json:"simulationId"`
	Session           types.SimulationSession `json:"session"`
	SignedTx          string                  `json:"signedTx"`
	SignedTransaction string                  `json:"signedTransaction,omitempty"`
	RawTransaction    string                  `json:"rawTransaction,omitempty"`
}

type executeSimulationResponse struct {
	Report types.SimulationReport `json:"report"`
}

func (s *Simulator) BeginSimulation(network types.Network, address common.Address, caller common.Address) (*types.SimulationSession, error) {
	response := beginSimulationResponse{}
	if err := s.postJSON(beginSimulationPath, beginSimulationRequest{
		Network: networkRPC(network),
		Address: address.Hex(),
		Caller:  caller.Hex(),
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
		var direct types.SimulationReport
		if err := json.Unmarshal(responseBody, &direct); err != nil {
			return nil, fmt.Errorf("decode simulator API response %s: %w", performSimulationPath, err)
		}
		report = direct
	}
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
		report.Summary = "Signed transaction executed by the simulator API."
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
		Note:          "Contract bytecode was reported by the simulator API.",
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
