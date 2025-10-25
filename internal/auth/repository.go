package auth

import (
	"context"
	"database/sql"
	"fmt"
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

type AuthRepository struct {
	storage Storage
}

func NewAuthRepository(storage Storage) *AuthRepository {
	return &AuthRepository{storage: storage}
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	op := "auth.repository.get_by_email"
	row := r.storage.QueryRowContext(ctx, "SELECT id, name, email, password, refresh_token FROM users WHERE email = $1 AND deleted_at IS NULL", email)
	user := &User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RefreshToken); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *AuthRepository) Create(ctx context.Context, user *User) error {
	op := "auth.repository.create"
	res, err := r.storage.ExecContext(
		ctx,
		"INSERT INTO users (name, email, password, refresh_token) VALUES($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RefreshToken,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rows, err := res.RowsAffected()
	if rows != 1 {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *AuthRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*User, error) {
	op := "auth.repository.getByRefreshToken"
	user := &User{}
	row := r.storage.QueryRowContext(ctx, "SELECT id, refresh_token FROM users WHERE refresh_token = $1 AND deleted_at IS NULL", refreshToken)
	if err := row.Scan(&user.ID, &user.RefreshToken); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, userID uint64, refreshToken string) error {
	op := "auth.repository.UpdateRefreshToken"
	res, err := r.storage.ExecContext(
		ctx,
		"UPDATE users SET refresh_token = $1 WHERE id = $2 AND deleted_at IS NULL",
		refreshToken, userID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rows, err := res.RowsAffected()
	if rows != 1 {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
