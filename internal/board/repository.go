package board

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrBoardAlreadyExists = errors.New("board already exists")
	ErrBoardNotFound      = errors.New("board not found")
	ErrColumnNotFound     = errors.New("column not found")
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

func (r *BoardRepository) GetList(
	ctx context.Context,
	filter *BoardGetFilter,
) (*BoardListResult, error) {
	op := "board.repository.GetList"
	boards := []*Board{}
	limit := filter.PerPage
	offset := (filter.Page - 1) * filter.PerPage
	query, params := buildQuery(filter)
	params = append(params, limit, offset)
	rows, err := r.storage.QueryContext(
		ctx,
		query,
		params...,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		board := &Board{}
		if err := rows.Scan(
			&board.ID,
			&board.Name,
			&board.Description,
			&board.CreatedAt,
			&board.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		boards = append(boards, board)
	}
	var count uint64
	countQuery := "SELECT COUNT(*) FROM boards WHERE user_id = $1 AND deleted_at IS NULL"
	where := []string{}
	params = []any{filter.UserID}
	where, params, _ = addFilterToQuery(where, params, filter.FilterFields)
	if len(where) > 0 {
		countQuery += " AND " + strings.Join(where, " AND ")
	}
	row := r.storage.QueryRowContext(ctx, countQuery, params...)
	err = row.Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	result := &BoardListResult{
		Data:       boards,
		TotalCount: count,
	}
	return result, nil
}

func buildQuery(filter *BoardGetFilter) (string, []any) {
	var where []string
	params := []any{filter.UserID}
	baseQuery := `
        SELECT id, name, description, created_at, updated_at
        FROM boards
        WHERE user_id = $1 AND deleted_at IS NULL
    `
	where, params, paramIndex := addFilterToQuery(where, params, filter.FilterFields)
	if len(where) > 0 {
		baseQuery += " AND " + strings.Join(where, " AND ")
	}
	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	return baseQuery, params
}

func addFilterToQuery(
	where []string,
	params []any,
	filter *Filters,
) ([]string, []any, int) {
	paramIndex := 2
	if filter != nil {
		if filter.Name != nil {
			where = append(where, fmt.Sprintf("name ILIKE $%d", paramIndex))
			params = append(params, "%"+*filter.Name+"%")
			paramIndex++
		}
		if filter.Description != nil {
			where = append(where, fmt.Sprintf("description ILIKE $%d", paramIndex))
			params = append(params, "%"+*filter.Description+"%")
			paramIndex++
		}
	}
	return where, params, paramIndex
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

func (r *BoardRepository) GetColumnList(ctx context.Context, uuid string) ([]*BoardColumn, error) {
	op := "board.repository.GetColumnList"
	columnsRaw := []*BoardColumn{}
	query := `
		SELECT id, board_id, name, color, position, created_at 
		FROM board_columns bc WHERE board_id = $1 AND deleted_at IS NULL
		ORDER BY position
	`
	rows, err := r.storage.QueryContext(ctx, query, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		column := &BoardColumn{}
		err := rows.Scan(
			&column.ID,
			&column.BoardID,
			&column.Name,
			&column.Color,
			&column.Position,
			&column.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		columnsRaw = append(columnsRaw, column)
	}
	return columnsRaw, nil
}

func (r *BoardRepository) CreateColumn(ctx context.Context, column *BoardColumn) error {
	op := "board.repository.CreateColumn"
	query := `
		INSERT INTO board_columns (board_id, name, color, position)
		VALUES ($1, $2, $3, $4)
	`
	result, err := r.storage.ExecContext(ctx, query, column.BoardID, column.Name, column.Color, column.Position)
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

func (r *BoardRepository) UpdateColumn(ctx context.Context, column *BoardColumn) error {
	op := "board.repository.UpdateColumn"
	query := `
		UPDATE board_columns
		SET
			name = $1, color = $2, position = $3, updated_at = NOW()
		WHERE
			board_id = $4 AND id = $5
			and deleted_at is NULL;
	`
	result, err := r.storage.ExecContext(
		ctx,
		query,
		column.Name,
		column.Color,
		column.Position,
		column.BoardID,
		column.ID,
	)
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

func (r *BoardRepository) DeleteColumn(ctx context.Context, column *BoardColumn) error {
	op := "board.repository.DeleteColumn"
	var exists bool
	row := r.storage.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM board_columns WHERE board_id = $1 AND id = $2)",
		column.BoardID,
		column.ID,
	)
	err := row.Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrColumnNotFound)
	}
	query := `
		UPDATE board_columns SET
			deleted_at = NOW()
		WHERE board_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	result, err := r.storage.ExecContext(ctx, query, column.BoardID, column.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("%s: %w", op, ErrColumnNotFound)
	}
	return nil
}

func (r *BoardRepository) MoveToColumn(ctx context.Context, id string, columnID, fromPosition, toPosition uint64) error {
	op := "board.repository.MoveToColumn"
	if fromPosition == toPosition {
		return nil
	}
	query := `
		UPDATE board_columns SET position = $1, updated_at = NOW() WHERE id = $2 AND board_id = $3 AND deleted_at IS NULL
	`
	result, err := r.storage.ExecContext(ctx, query, toPosition, columnID, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("%s: %w", op, ErrColumnNotFound)
	}
	return nil
}
