package testmocks

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE accounts (
	id TEXT PRIMARY KEY,
	address TEXT NOT NULL UNIQUE,
	last_used INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE networks (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	rpc TEXT NOT NULL,
	symbol TEXT NOT NULL,
	chain_id TEXT NOT NULL,
	block_explorer TEXT
);

CREATE TABLE token_transactions (
	id TEXT PRIMARY KEY,
	tx_hash TEXT NOT NULL UNIQUE,
	recipient TEXT NOT NULL,
	amount TEXT NOT NULL,
	account_id TEXT NOT NULL,
	network_id TEXT NOT NULL,

	FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE
);

CREATE TABLE contact_transactions (
	id TEXT PRIMARY KEY,
	tx_hash TEXT NOT NULL UNIQUE,
	recipient TEXT NOT NULL,
	token TEXT NULL,
	amount TEXT NOT NULL,
	account_id TEXT NOT NULL,
	network_id TEXT NOT NULL,

	FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE
);
`

type Storage struct {
	db *sql.DB
}

func NewStorage(t *testing.T) *Storage {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() returned error: %v", err)
	}
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("mock storage schema setup returned error: %v", err)
	}

	storage := &Storage{db: db}
	t.Cleanup(func() {
		if err := storage.Close(); err != nil {
			t.Fatalf("mock storage Close() returned error: %v", err)
		}
	})

	return storage
}

func (s *Storage) Init() error {
	return nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Storage) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *Storage) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *Storage) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *Storage) DB() *sql.DB {
	return s.db
}
