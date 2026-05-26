PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS accounts (
	id TEXT PRIMARY KEY,
	address TEXT NOT NULL UNIQUE,
	last_used INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS networks (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	rpc TEXT NOT NULL,
	symbol TEXT NOT NULL,
	chain_id TEXT NOT NULL,
	block_explorer INTEGER
);

CREATE TABLE IF NOT EXISTS tokens (
	id TEXT PRIMARY KEY,
	symbol TEXT NOT NULL,
	name TEXT NOT NULL,
	address TEXT NOT NULL,
	network TEXT NOT NULL,

	FOREIGN KEY (network) REFERENCES networks(id) ON DELETE CASCADE,
	UNIQUE(address, network)
);

CREATE TABLE IF NOT EXISTS token_transactions (
	id TEXT PRIMARY KEY,
	tx_hash TEXT NOT NULL UNIQUE,
	amount TEXT NOT NULL,
	account_id TEXT NOT NULL,
	network_id TEXT NOT NULL,

	FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS markets (
	id TEXT PRIMARY KEY,
	network TEXT NOT NULL,
	router TEXT NOT NULL,

	FOREIGN KEY (network) REFERENCES networks(id) ON DELETE CASCADE,
	UNIQUE(network, router)
);

CREATE TABLE IF NOT EXISTS contacts (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS contact_transactions (
	id TEXT PRIMARY KEY,
	tx_hash TEXT NOT NULL UNIQUE,
	token TEXT NULL,
	amount TEXT NOT NULL,
	account_id TEXT NOT NULL,
	network_id TEXT NOT NULL,

	FOREIGN KEY (token) REFERENCES tokens(id) ON DELETE SET NULL,
	FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tokens_network
ON tokens(network);

CREATE INDEX IF NOT EXISTS idx_token_transactions_account
ON token_transactions(account_id);

CREATE INDEX IF NOT EXISTS idx_token_transactions_network
ON token_transactions(network_id);

CREATE INDEX IF NOT EXISTS idx_markets_network
ON markets(network);

CREATE INDEX IF NOT EXISTS idx_contact_transactions_account
ON contact_transactions(account_id);

CREATE INDEX IF NOT EXISTS idx_contact_transactions_network
ON contact_transactions(network_id);

CREATE INDEX IF NOT EXISTS idx_contact_transactions_token
ON contact_transactions(token);
