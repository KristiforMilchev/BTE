package repositories_test

import (
	"context"
	"testing"

	"bos/repositories"
	"bos/tests/testmocks"
	"bos/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
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
		"INSERT INTO networks (id, name, rpc, symbol, chain_id, block_explorer) VALUES (?, ?, ?, ?, ?, ?)",
		"net-1",
		"Local",
		rpc,
		"ETH",
		1337,
		"http://explorer",
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

func TestTransactionsRepositorySaveNativeTransferStoresAccountAndNetwork(t *testing.T) {
	storage := testmocks.NewStorage(t)
	register := repositories.NewRegister(storage)
	account := common.HexToAddress("0x1111111111111111111111111111111111111111")
	networkID := uuid.New()
	rpc := "http://localhost:8545"

	if err := register.Accounts.AddAccount(account); err != nil {
		t.Fatalf("AddAccount() returned error: %v", err)
	}
	_, err := storage.Exec(
		context.Background(),
		"INSERT INTO networks (id, name, rpc, symbol, chain_id, block_explorer) VALUES (?, ?, ?, ?, ?, ?)",
		networkID.String(),
		"Local",
		rpc,
		"ETH",
		1337,
		"http://explorer",
	)
	if err != nil {
		t.Fatalf("storage.Exec() returned error: %v", err)
	}

	err = register.Transactions.SaveNativeTransfer(
		types.TxDraft{
			FromAddress:      account.Hex(),
			RecipientAddress: "0x2222222222222222222222222222222222222222",
			Amount:           "1.25",
		},
		"0xtx",
		types.Network{Id: networkID, Rpc: &rpc},
	)
	if err != nil {
		t.Fatalf("SaveNativeTransfer() returned error: %v", err)
	}

	var txHash string
	var recipient string
	var amount string
	var savedAccountID string
	var savedNetworkID string
	err = storage.QueryRow(context.Background(), `
		SELECT tx_hash, recipient, amount, account_id, network_id
		FROM contact_transactions
		WHERE tx_hash = ?;
	`, "0xtx").Scan(&txHash, &recipient, &amount, &savedAccountID, &savedNetworkID)
	if err != nil {
		t.Fatalf("QueryRow() returned error: %v", err)
	}
	if txHash != "0xtx" || amount != "1.25" {
		t.Fatalf("saved transaction = (%q, %q), want (0xtx, 1.25)", txHash, amount)
	}
	if recipient != "0x2222222222222222222222222222222222222222" {
		t.Fatalf("recipient = %q, want 0x2222222222222222222222222222222222222222", recipient)
	}
	if savedNetworkID != networkID.String() {
		t.Fatalf("network_id = %q, want %q", savedNetworkID, networkID.String())
	}

	var accountID string
	err = storage.QueryRow(context.Background(), "SELECT id FROM accounts WHERE address = ?;", account.Hex()).Scan(&accountID)
	if err != nil {
		t.Fatalf("account lookup returned error: %v", err)
	}
	if savedAccountID != accountID {
		t.Fatalf("account_id = %q, want %q", savedAccountID, accountID)
	}
}

func TestTransactionsRepositoryGetTransactionsReturnsTokenAndContactTransactions(t *testing.T) {
	storage := testmocks.NewStorage(t)
	register := repositories.NewRegister(storage)
	account := common.HexToAddress("0x1111111111111111111111111111111111111111")
	networkID := uuid.New()

	if err := register.Accounts.AddAccount(account); err != nil {
		t.Fatalf("AddAccount() returned error: %v", err)
	}
	_, err := storage.Exec(
		context.Background(),
		"INSERT INTO networks (id, name, rpc, symbol, chain_id, block_explorer) VALUES (?, ?, ?, ?, ?, ?)",
		networkID.String(),
		"Local",
		"http://localhost:8545",
		"ETH",
		1337,
		"http://explorer",
	)
	if err != nil {
		t.Fatalf("insert network returned error: %v", err)
	}

	var accountID string
	err = storage.QueryRow(context.Background(), "SELECT id FROM accounts WHERE address = ?;", account.Hex()).Scan(&accountID)
	if err != nil {
		t.Fatalf("account lookup returned error: %v", err)
	}

	_, err = storage.Exec(
		context.Background(),
		"INSERT INTO token_transactions (id, tx_hash, recipient, amount, account_id, network_id) VALUES (?, ?, ?, ?, ?, ?)",
		uuid.NewString(),
		"0xtoken",
		"0x2222222222222222222222222222222222222222",
		"5 TOKEN",
		accountID,
		networkID.String(),
	)
	if err != nil {
		t.Fatalf("insert token transaction returned error: %v", err)
	}
	_, err = storage.Exec(
		context.Background(),
		"INSERT INTO contact_transactions (id, tx_hash, recipient, token, amount, account_id, network_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		uuid.NewString(),
		"0xcontact",
		"0x3333333333333333333333333333333333333333",
		nil,
		"1 ETH",
		accountID,
		networkID.String(),
	)
	if err != nil {
		t.Fatalf("insert contact transaction returned error: %v", err)
	}

	accountHex := account.Hex()
	got, err := register.Transactions.GetTransactions(&networkID, &accountHex)
	if err != nil {
		t.Fatalf("GetTransactions() returned error: %v", err)
	}
	if len(*got) != 2 {
		t.Fatalf("GetTransactions() returned %d transactions, want 2", len(*got))
	}
	if (*got)[0].To != "0x2222222222222222222222222222222222222222" || (*got)[0].TxHash != "0xtoken" || (*got)[0].Amount != "5 TOKEN" {
		t.Fatalf("first transaction = %+v, want token transaction", (*got)[0])
	}
	if (*got)[1].To != "0x3333333333333333333333333333333333333333" || (*got)[1].TxHash != "0xcontact" || (*got)[1].Amount != "1 ETH" {
		t.Fatalf("second transaction = %+v, want contact transaction", (*got)[1])
	}
}
