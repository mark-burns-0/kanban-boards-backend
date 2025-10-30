package comment

import (
	"backend/internal/shared/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrCardNotFound         = errors.New("card not found")
	ErrCommentAlreadyExists = errors.New("comment already exists")
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

type CommentRepository struct {
	storage Storage
}

func NewCommentRepository(
	storage Storage,
) *CommentRepository {
	return &CommentRepository{
		storage: storage,
	}
}

func (r *CommentRepository) Create(ctx context.Context, comment *Comment) error {
	const op = "comment.repository.Create"
	query := `INSERT INTO comments (card_id, user_id, text) VALUES($1, $2,$3)`
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		ErrCommentAlreadyExists,
		comment.CardID,
		comment.UserID,
		comment.Text,
	)
}
func (r *CommentRepository) Update(ctx context.Context, comment *Comment) error {
	const op = "comment.repository.Update"
	query := "UPDATE comments SET text = $1, updated_at = NOW() WHERE id = $2 AND card_id = $3 AND user_id = $4"
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		ErrCardNotFound,
		comment.Text,
		comment.ID,
		comment.CardID,
		comment.UserID,
	)
}

func (r *CommentRepository) Delete(ctx context.Context, commentID uint64) error {
	const op = "comment.repository.Delete"
	existQuery := "SELECT EXISTS (SELECT 1 FROM comments WHERE id = $1)"
	row := r.storage.QueryRowContext(ctx, existQuery, commentID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrCardNotFound)
	}
	query := "UPDATE comments SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"
	_, err := r.storage.ExecContext(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
