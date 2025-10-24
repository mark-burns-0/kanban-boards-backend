package board

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrBoardAlreadyExists = errors.New("board already exists")
	ErrBoardNotFound      = errors.New("board not found")
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

type BoardRepository struct {
	storage Storage
}

func NewBoardRepository(storage Storage) *BoardRepository {
	return &BoardRepository{
		storage: storage,
	}
}

func (r *BoardRepository) Get(ctx context.Context, uuid string) (*Board, error) {
	return nil, nil
}

func (r *BoardRepository) GetList(ctx context.Context, userID uint64) ([]*Board, error) {
	return nil, nil
}

func (r *BoardRepository) Create(ctx context.Context, board *Board) error {
	op := "board.repository.Create"
	query := "INSERT INTO boards (name, description, user_id) VALUES ($1, $2, $3)"
	result, err := r.storage.ExecContext(ctx, query, board.Name, board.Description, board.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("%s: %w", op, ErrBoardAlreadyExists)
	}
	return nil
}

func (r *BoardRepository) Update(ctx context.Context, board *Board) error {
	op := "board.repository.Update"
	query := "UPDATE boards SET name = $1, description = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL"
	fmt.Println(board)
	result, err := r.storage.ExecContext(ctx, query, board.Name, board.Description, board.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *BoardRepository) Delete(ctx context.Context, uuid string) error {
	op := "board.repository.Delete"
	var exists bool
	row := r.storage.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM boards WHERE id = $1)",
		uuid,
	)
	err := row.Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrBoardNotFound)
	}
	query := "UPDATE boards SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"
	_, err = r.storage.ExecContext(ctx, query, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *BoardRepository) MoveToColumn(ctx context.Context, id string, columnID, fromPosition, toPosition uint64) error {
	op := "board.repository.MoveToColumn"
	fmt.Println(op)
	return nil
}
