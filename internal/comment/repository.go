package comment

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

func (r *CommentRepository) Create() {}
func (r *CommentRepository) Update() {}
func (r *CommentRepository) Delete() {}
