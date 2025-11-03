package repository

import (
	"backend/internal/user/domain"
	"context"
	"database/sql"
	"fmt"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type UserRepository struct {
	storage                Storage
	updateStmt             *sql.Stmt
	updateWithPasswordStmt *sql.Stmt
}

func NewUserRepository(storage Storage) (*UserRepository, error) {
	const op = "user.repository.NewUserRepository"
	var err error
	repo := &UserRepository{
		storage: storage,
	}

	repo.updateStmt, err = storage.GetDB().Prepare(`
		UPDATE users SET name = $1, email = $2, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	repo.updateWithPasswordStmt, err = storage.GetDB().Prepare(`
		UPDATE users SET name = $1, email = $2, password = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return repo, nil
}

func (r *UserRepository) Get(ctx context.Context, id uint64) (*domain.User, error) {
	const op = "user.repository.Get"
	user := &domain.User{}

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

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	const op = "user.repository.Update"

	if user.Password != "" {
		res, err := r.updateWithPasswordStmt.ExecContext(
			ctx, user.Name, user.Email, user.Password, user.ID,
		)
		return handleQueryExec(op, res, err)
	}

	res, err := r.updateStmt.ExecContext(
		ctx, user.Name, user.Email, user.ID,
	)

	return handleQueryExec(op, res, err)
}

func handleQueryExec(op string, res sql.Result, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	affected, err := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}

	return nil
}
