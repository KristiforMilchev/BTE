package repositories_test

import (
	"context"
	"testing"

	"bos/repositories"
	"bos/tests/testmocks"

	"github.com/ethereum/go-ethereum/common"
)

func TestAccountsRepositoryAccountReturnsErrorWhenEmpty(t *testing.T) {
	register := repositories.NewRegister(testmocks.NewStorage(t))

	if got, err := register.Accounts.Account(); err == nil {
		t.Fatalf("Account() = %v, want error", got)
	}
}

func TestAccountsRepositoryAddAccountTracksLastUsed(t *testing.T) {
	register := repositories.NewRegister(testmocks.NewStorage(t))
	first := common.HexToAddress("0x1111111111111111111111111111111111111111")
	second := common.HexToAddress("0x2222222222222222222222222222222222222222")

	if err := register.Accounts.AddAccount(first); err != nil {
		t.Fatalf("AddAccount(first) returned error: %v", err)
	}
	got, err := register.Accounts.Account()
	if err != nil {
		t.Fatalf("Account() returned error after first add: %v", err)
	}
	if *got != first {
		t.Fatalf("Account() = %s, want %s", got.Hex(), first.Hex())
	}

	if err := register.Accounts.AddAccount(second); err != nil {
		t.Fatalf("AddAccount(second) returned error: %v", err)
	}
	got, err = register.Accounts.Account()
	if err != nil {
		t.Fatalf("Account() returned error after second add: %v", err)
	}
	if *got != second {
		t.Fatalf("Account() = %s, want %s", got.Hex(), second.Hex())
	}
}

func TestNetworkRepositoryNameReturnsNetworkNameByRPC(t *testing.T) {
	storage := testmocks.NewStorage(t)
	rpc := "http://localhost:8545"
	_, err := storage.Exec(
		context.Background(),
		"INSERT INTO networks (id, name, rpc) VALUES (?, ?, ?)",
		"net-1",
		"Local",
		rpc,
	)
	if err != nil {
		t.Fatalf("storage.Exec() returned error: %v", err)
	}

	repo := repositories.NewNetworkRepository(storage)
	got, err := repo.Name(&rpc)
	if err != nil {
		t.Fatalf("Name() returned error: %v", err)
	}
	if *got != "Local" {
		t.Fatalf("Name() = %q, want %q", *got, "Local")
	}
}
