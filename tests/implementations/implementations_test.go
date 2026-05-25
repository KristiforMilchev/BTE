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

			return jsonResponse(`{"simulationId":"session-1","clonedRpc":"http://sim-rpc","transaction":{"chainId":"0x539"},"contract":{"hasCode":true},"runtimeHex":"0x6001","runtimeBinary":"0b0110000000000001"}`), nil

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

			return jsonResponse(`{"report":{"title":"API Report","status":"ok","bytecodeChecks":[{"address":"0x2222222222222222222222222222222222222222","isContract":true,"runtimeHex":"0x6001","runtimeBinary":"0b0110000000000001"}]}}`), nil

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

	report, err := simulator.ExecuteSignedTransaction(*session, []byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Fatalf("ExecuteSignedTransaction() returned error: %v", err)
	}
	if report.Title != "API Report" || len(report.BytecodeChecks) != 1 {
		t.Fatalf("report = %+v", report)
	}
	if !beginCalled || !executeCalled {
		t.Fatalf("beginCalled=%v executeCalled=%v", beginCalled, executeCalled)
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
