package implementations

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Storage struct {
	path       string
	schemaPath string
	db         *sql.DB
}

func NewStorage(path string, schemaPath string) *Storage {
	return &Storage{
		path:       path,
		schemaPath: schemaPath,
	}
}

func (s *Storage) Init() error {
	if s.path == "" {
		s.path = "./bos.db"
	}

	dir := filepath.Dir(s.path)

	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	db, err := sql.Open("sqlite", s.path)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	schema, err := os.ReadFile(s.schemaPath)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(string(schema))
	if err != nil {
		return err
	}

	return s.migrate()
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

func (s *Storage) migrate() error {
	if err := s.ensureColumn("token_transactions", "network_id", "TEXT"); err != nil {
		return err
	}
	if err := s.ensureColumn("contact_transactions", "network_id", "TEXT"); err != nil {
		return err
	}
	if err := s.ensureColumn("token_transactions", "recipient", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := s.ensureColumn("contact_transactions", "recipient", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}

	_, err := s.db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_token_transactions_network
		ON token_transactions(network_id);

		CREATE INDEX IF NOT EXISTS idx_contact_transactions_network
		ON contact_transactions(network_id);
	`)
	return err
}

func (s *Storage) ensureColumn(table string, column string, definition string) error {
	rows, err := s.db.Query("PRAGMA table_info(" + table + ");")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var columnType string
		var notNull int
		var defaultValue any
		var primaryKey int
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return err
		}
		if name == column {
			return rows.Err()
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	_, err = s.db.Exec("ALTER TABLE " + table + " ADD COLUMN " + column + " " + definition + ";")
	return err
}
