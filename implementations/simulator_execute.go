package implementations

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"bos/dtos"
	"bos/types"
)

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
