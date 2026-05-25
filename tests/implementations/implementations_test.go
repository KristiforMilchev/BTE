package implementations_test

import (
	"context"
	"path/filepath"
	"testing"

	"bos/implementations"
	"bos/repositories"

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
