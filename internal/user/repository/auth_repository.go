package repository

import (
	"backend/internal/shared/ports/repository"
	"backend/internal/shared/utils"
	"backend/internal/user/domain"
	"context"
	"fmt"
)

type AuthRepository struct {
	storage repository.Storage
}

func NewAuthRepository(storage repository.Storage) *AuthRepository {
	return &AuthRepository{storage: storage}
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op = "auth.repository.get_by_email"
	row := r.storage.QueryRowContext(ctx, "SELECT id, name, email, password, refresh_token FROM users WHERE email = $1 AND deleted_at IS NULL", email)
	user := &domain.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RefreshToken); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *AuthRepository) Create(ctx context.Context, user *domain.User) error {
	const op = "auth.repository.create"
	query := "INSERT INTO users (name, email, password, refresh_token) VALUES($1, $2, $3, $4)"
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		domain.ErrUserAlreadyExists,
		user.Name,
		user.Email,
		user.Password,
		user.RefreshToken,
	)
}

func (r *AuthRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error) {
	const op = "auth.repository.getByRefreshToken"
	user := &domain.User{}
	row := r.storage.QueryRowContext(ctx, "SELECT id, refresh_token FROM users WHERE refresh_token = $1 AND deleted_at IS NULL", refreshToken)
	if err := row.Scan(&user.ID, &user.RefreshToken); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, userID uint64, refreshToken string) error {
	const op = "auth.repository.UpdateRefreshToken"
	query := "UPDATE users SET refresh_token = $1 WHERE id = $2 AND deleted_at IS NULL"
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		domain.ErrUserNotFound,
		refreshToken,
		userID,
	)
}
