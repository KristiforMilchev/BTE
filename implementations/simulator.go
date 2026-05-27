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
	"bos/dtos"
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

func (s *Simulator) ExecuteSignedTransaction(session types.SimulationSession, signedTx []byte) (*types.SimulationReport, error) {
	if len(signedTx) == 0 {
		return nil, fmt.Errorf("signed transaction is empty")
	}

	rawTx := "0x" + hex.EncodeToString(signedTx)
	responseBody, err := s.postJSONRaw(performSimulationPath, dtos.SimulatorExecuteRequest{
		SimulationID:      session.ID,
		Session:           session,
		SignedTx:          rawTx,
		SignedTransaction: rawTx,
		RawTransaction:    rawTx,
	})
	if err != nil {
		return nil, err
	}

	var wrapped dtos.SimulatorExecuteResponse
	_ = json.Unmarshal(responseBody, &wrapped)
	report := wrapped.Report
	if report.IsEmpty() {
		var apiResponse dtos.SimulatorPerformResponse
		if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
			return nil, fmt.Errorf("decode simulator API response %s: %w", performSimulationPath, err)
		}
		if apiResponse.Execution.Status != "" || apiResponse.SimulationID != "" {
			report = apiResponse.ToSimulationReport(session)
		} else {
			var direct types.SimulationReport
			if err := json.Unmarshal(responseBody, &direct); err != nil {
				return nil, fmt.Errorf("decode simulator API response %s: %w", performSimulationPath, err)
			}
			report = direct
		}
	}
	report = report.WithFallbacks(session)
	return &report, nil
}

func (s *Simulator) ExecuteSignedTransactions(session types.SimulationSession, signedTxs []types.SignedSimulationTransaction) (*types.SimulationReport, error) {
	if len(signedTxs) == 0 {
		return nil, fmt.Errorf("signed transactions are empty")
	}

	responseBody, err := s.postJSONRaw(performSimulationPath, dtos.SimulatorExecuteRequest{
		SimulationID:       session.ID,
		Session:            session,
		SignedTransactions: signedTxs,
	})
	if err != nil {
		return nil, err
	}

	var apiResponse dtos.SimulatorPerformResponse
	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return nil, fmt.Errorf("decode simulator API response %s: %w", performSimulationPath, err)
	}

	report := apiResponse.ToSimulationReport(session).WithFallbacks(session)
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
