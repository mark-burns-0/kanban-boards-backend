package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

type StorageConfig interface {
	GetUsername() string
	GetPassword() string
	GetHost() string
	GetPort() string
	GetName() string
	GetSSLMode() string
}

func NewStorage(cfg StorageConfig) (*Storage, error) {
	const op = "storage.NewStorage"
	dsn := buildDSN(cfg)
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Exec(query string, args ...any) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

func (s *Storage) Query(query string, args ...any) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

func (s *Storage) QueryRow(query string, args ...any) *sql.Row {
	return s.db.QueryRow(query, args...)
}

func (s *Storage) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}

func (s *Storage) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *Storage) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *Storage) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *Storage) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, opts)
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}

func buildDSN(cfg StorageConfig) string {
	dsn := strings.Builder{}
	dsn.WriteString("postgres://")
	dsn.WriteString(cfg.GetUsername())
	dsn.WriteString(":")
	dsn.WriteString(cfg.GetPassword())
	dsn.WriteString("@")
	dsn.WriteString(cfg.GetHost())
	dsn.WriteString(":")
	dsn.WriteString(cfg.GetPort())
	dsn.WriteString("/")
	dsn.WriteString(cfg.GetName())
	dsn.WriteString("?sslmode=")
	dsn.WriteString(cfg.GetSSLMode())
	return dsn.String()
}
