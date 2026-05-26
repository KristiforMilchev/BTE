package interfaces

import (
	"context"
	"database/sql"
)

type IStorage interface {
	Init() error
	Close() error

	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row

	DB() *sql.DB
}
