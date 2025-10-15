package user

import (
	"context"
	"database/sql"
)

type Storage interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)

	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type UserRepository struct {
	storage Storage
}

func NewUserRepository(storage Storage) *UserRepository {
	return &UserRepository{
		storage: storage,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*User, error) {
	return nil, nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	return nil
}
