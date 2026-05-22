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
