package implementations_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"path/filepath"
	"testing"

	"bos/implementations"
	"bos/repositories"
	"bos/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

func TestStorageInitCreatesSchemaAndExecutesQueries(t *testing.T) {
	storage := implementations.NewStorage(
		filepath.Join(t.TempDir(), "bos.db"),
		filepath.Join("..", "..", "sql", "schema.sql"),
	)

	if err := storage.Init(); err != nil {
		t.Fatalf("Init() returned error: %v", err)
	}
	t.Cleanup(func() {
		if err := storage.Close(); err != nil {
			t.Fatalf("Close() returned error: %v", err)
		}
	})

	_, err := storage.Exec(
		context.Background(),
		"INSERT INTO networks (id, name, rpc, symbol, chain_id, block_explorer) VALUES (?, ?, ?, ?, ?, ?)",
		"net-1",
		"Local",
		"http://localhost:8545",
		"ETH",
		1337,
		"http://explorer",
	)
	if err != nil {
		t.Fatalf("Exec() returned error: %v", err)
	}

	var name string
	err = storage.QueryRow(context.Background(), "SELECT name FROM networks WHERE id = ?", "net-1").Scan(&name)
	if err != nil {
		t.Fatalf("QueryRow().Scan() returned error: %v", err)
	}
	if name != "Local" {
		t.Fatalf("stored network name = %q, want %q", name, "Local")
	}
}

func TestStorageCloseBeforeInitIsNoop(t *testing.T) {
	storage := implementations.NewStorage("", "")
	if err := storage.Close(); err != nil {
		t.Fatalf("Close() before Init() returned error: %v", err)
	}
}

func TestNetworkProviderStartsWithFirstSavedNetwork(t *testing.T) {
	storage := implementations.NewStorage(
		filepath.Join(t.TempDir(), "bos.db"),
		filepath.Join("..", "..", "sql", "schema.sql"),
	)
	if err := storage.Init(); err != nil {
		t.Fatalf("Init() returned error: %v", err)
	}
	t.Cleanup(func() {
		if err := storage.Close(); err != nil {
			t.Fatalf("Close() returned error: %v", err)
		}
	})

	_, err := storage.Exec(
		context.Background(),
		"INSERT INTO networks (id, name, rpc, symbol, chain_id, block_explorer) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.NewString(),
		"First",
		"http://first",
		"ETH",
		1,
		"http://explorer",
	)
	if err != nil {
		t.Fatalf("Exec(first network) returned error: %v", err)
	}

	provider := implementations.NewNetworkProvider(repositories.NewNetworkRepository(storage))
	active := provider.Network()

	if active.Name == nil || *active.Name != "First" {
		t.Fatalf("active network name = %v, want First", active.Name)
	}
	if active.Rpc == nil || *active.Rpc != "http://first" {
		t.Fatalf("active network rpc = %v, want http://first", active.Rpc)
	}
	if active.Chain == nil || active.Chain.Int64() != 1 {
		t.Fatalf("active network chain = %v, want 1", active.Chain)
	}
}

func TestERC20Address(t *testing.T) {
	address := common.HexToAddress("0x1111111111111111111111111111111111111111")
	token, err := implementations.NewERC20(nil, address)
	if err != nil {
		t.Fatalf("NewERC20() returned error: %v", err)
	}
	if got := token.Address(); got != address {
		t.Fatalf("Address() = %s, want %s", got.Hex(), address.Hex())
	}
}

func TestUniswapV2RouterAddress(t *testing.T) {
	address := common.HexToAddress("0x2222222222222222222222222222222222222222")
	router, err := implementations.NewUniswapV2Router(nil, address)
	if err != nil {
		t.Fatalf("NewUniswapV2Router() returned error: %v", err)
	}
	if got := router.Address(); got != address {
		t.Fatalf("Address() = %s, want %s", got.Hex(), address.Hex())
	}
}

func TestSimulatorCallsBeginAndExecuteAPI(t *testing.T) {
	var beginCalled bool
	var executeCalled bool

	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/v1/simulation/begin":
			beginCalled = true
			if r.Method != http.MethodPost {
				t.Fatalf("begin method = %s, want POST", r.Method)
			}

			var request struct {
				Address string `json:"address"`
				Caller  string `json:"caller"`
				Network string `json:"network"`
			}
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatalf("decode begin request: %v", err)
			}
			if request.Address != "0x2222222222222222222222222222222222222222" {
				t.Fatalf("begin address = %q", request.Address)
			}
			if request.Caller != "0x1111111111111111111111111111111111111111" {
				t.Fatalf("begin caller = %q", request.Caller)
			}
			if request.Network != "http://source-rpc" {
				t.Fatalf("begin network = %q", request.Network)
			}

			return jsonResponse(`{"simulationId":"session-1","clonedRpc":"http://sim-rpc","transaction":{"type":"legacy","chainId":"0x539","from":"0x1111111111111111111111111111111111111111","to":"0x2222222222222222222222222222222222222222","value":"0x0","data":"0x","nonce":"0x7","gas":"0x5208","gasPrice":"0x3b9aca00"},"contract":{"hasCode":true},"runtimeHex":"0x6001","runtimeBinary":"0b0110000000000001"}`), nil

		case "/v1/simulation/perform":
			executeCalled = true
			if r.Method != http.MethodPost {
				t.Fatalf("execute method = %s, want POST", r.Method)
			}

			var request struct {
				SimulationID      string                  `json:"simulationId"`
				Session           types.SimulationSession `json:"session"`
				SignedTx          string                  `json:"signedTx"`
				SignedTransaction string                  `json:"signedTransaction"`
				RawTransaction    string                  `json:"rawTransaction"`
			}
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatalf("decode execute request: %v", err)
			}
			if request.Session.ID != "session-1" || request.Session.RPC != "http://sim-rpc" {
				t.Fatalf("execute session = %+v", request.Session)
			}
			if request.SimulationID != "session-1" {
				t.Fatalf("simulationId = %q, want session-1", request.SimulationID)
			}
			if request.SignedTx != "0x010203" {
				t.Fatalf("signed tx = %q, want 0x010203", request.SignedTx)
			}
			if request.SignedTransaction != "0x010203" || request.RawTransaction != "0x010203" {
				t.Fatalf("signed tx aliases = %q %q, want 0x010203", request.SignedTransaction, request.RawTransaction)
			}

			return jsonResponse(`{"simulationId":"session-1","network":"http://source-rpc","rawTransactionSha256":"0xabc","balances":{"callerBefore":"0x10","callerAfter":"0x09","addressBefore":"0x00","addressAfter":"0x00"},"approvalFindings":[{"type":"erc20-approval","selector":"0x095ea7b3","description":"approval selector found in signed transaction","severity":"high"}],"contract":{"address":"0x2222222222222222222222222222222222222222","hasCode":true,"codeHashSha256":"0xcode","stateChanges":["receipt.status=0x1"]},"execution":{"mode":"ganache-fork","broadcasted":false,"transactionHash":"0xtx","status":"0x1","details":"transaction executed on Ganache fork; gas used 0x5208","forkBackendNeeded":false},"warnings":["signed transaction was executed only on the Ganache fork and was not broadcast to the upstream network"]}`), nil

		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		return nil, nil
	})

	name := "Local"
	rpc := "http://source-rpc"
	symbol := "ETH"
	network := types.Network{
		Name:   &name,
		Rpc:    &rpc,
		Symbol: &symbol,
		Chain:  big.NewInt(1337),
	}

	simulator := implementations.NewSimulatorWithHTTPClient("http://simulator.test", &http.Client{Transport: transport})
	session, err := simulator.BeginSimulation(
		network,
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
	)
	if err != nil {
		t.Fatalf("BeginSimulation() returned error: %v", err)
	}
	if !session.AddressContract || session.RuntimeHex != "0x6001" || session.ChainID != "0x539" {
		t.Fatalf("session bytecode fields = %+v", session)
	}
	if session.Transaction.To != "0x2222222222222222222222222222222222222222" || session.Transaction.Nonce != "0x7" {
		t.Fatalf("session transaction = %+v", session.Transaction)
	}

	report, err := simulator.ExecuteSignedTransaction(*session, []byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Fatalf("ExecuteSignedTransaction() returned error: %v", err)
	}
	if report.Title != "Contract Simulation" || report.RiskLevel != "High" || len(report.BytecodeChecks) != 1 {
		t.Fatalf("report = %+v", report)
	}
	if len(report.TokenApprovals) != 1 || report.TokenApprovals[0].Spender != "0x095ea7b3" {
		t.Fatalf("token approvals = %+v", report.TokenApprovals)
	}
	if len(report.Events) == 0 || report.Events[0].Details != "receipt.status=0x1" {
		t.Fatalf("report = %+v", report)
	}
	if !beginCalled || !executeCalled {
		t.Fatalf("beginCalled=%v executeCalled=%v", beginCalled, executeCalled)
	}
}

func TestSimulatorNativeTransferReportIsNotContractSimulation(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/v1/simulation/perform" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		return jsonResponse(`{"simulationId":"session-1","network":"http://source-rpc","rawTransactionSha256":"0xabc","balances":{"callerBefore":"0x10","callerAfter":"0x09","addressBefore":"0x00","addressAfter":"0x01"},"contract":{"address":"0x2222222222222222222222222222222222222222","hasCode":false},"execution":{"mode":"ganache-fork","broadcasted":false,"transactionHash":"0xtx","status":"0x1","forkBackendNeeded":false}}`), nil
	})

	name := "Local"
	rpc := "http://source-rpc"
	symbol := "ETH"
	session := types.SimulationSession{
		ID:      "session-1",
		Address: "0x2222222222222222222222222222222222222222",
		Caller:  "0x1111111111111111111111111111111111111111",
		Amount:  "0.1",
		Asset:   "ETH",
		Network: types.Network{
			Name:   &name,
			Rpc:    &rpc,
			Symbol: &symbol,
			Chain:  big.NewInt(1337),
		},
		Transaction: types.LedgerTransaction{
			From:  "0x1111111111111111111111111111111111111111",
			To:    "0x2222222222222222222222222222222222222222",
			Value: "0x16345785d8a0000",
		},
	}

	simulator := implementations.NewSimulatorWithHTTPClient("http://simulator.test", &http.Client{Transport: transport})
	report, err := simulator.ExecuteSignedTransaction(session, []byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Fatalf("ExecuteSignedTransaction() returned error: %v", err)
	}

	if report.Title != "Transaction Simulation" {
		t.Fatalf("report title = %q, want Transaction Simulation", report.Title)
	}
	if len(report.Calls) != 1 || report.Calls[0].Function != "regular wallet transfer" {
		t.Fatalf("report calls = %+v", report.Calls)
	}
	if len(report.BalanceChanges) != 2 || report.BalanceChanges[1].Asset != "Recipient native balance" {
		t.Fatalf("balance changes = %+v", report.BalanceChanges)
	}
	if len(report.BytecodeChecks) != 1 || report.BytecodeChecks[0].IsContract {
		t.Fatalf("bytecode checks = %+v", report.BytecodeChecks)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
