package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

var (
	ErrNotFound = fmt.Errorf("not found")
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

func (r *UserRepository) Get(ctx context.Context, id uint64) (*User, error) {
	op := "user.repository.Get"
	user := &User{}
	row := r.storage.QueryRowContext(
		ctx,
		"SELECT id, email, name, created_at FROM users WHERE id = $1 AND deleted_at IS NULL",
		id,
	)

	if err := row.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err.Error())
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	op := "user.repository.Update"
	query := strings.Builder{}
	requiredFields := 2
	args := []any{}
	args = append(args, user.Name, user.Email)

	query.WriteString("UPDATE users SET name = $1, email = $2, updated_at = NOW()")
	if user.Password != "" {
		requiredFields++
		query.WriteString(fmt.Sprintf(", password = $%d", requiredFields))
		args = append(args, user.Password)
	}
	args = append(args, user.ID)
	requiredFields++
	query.WriteString(fmt.Sprintf(" WHERE id = $%d", requiredFields))
	query.WriteString("AND deleted_at IS NULL")

	res, err := r.storage.ExecContext(
		ctx,
		query.String(),
		args...,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	affected, err := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}

	return nil
}
