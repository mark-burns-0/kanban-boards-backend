package card

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

type CardRepository struct {
	storage Storage
}

func NewCardRepository(
	storage Storage,
) *CardRepository {
	return &CardRepository{
		storage: storage,
	}
}

func (r *CardRepository) Create() {}

func (r *CardRepository) GetList() {}

func (r *CardRepository) Delete(id uint64) {}

func (r *CardRepository) MoveToNewPosition(
	cardID, boardID, toColumnID, cardPosition uint64,
) {
}
