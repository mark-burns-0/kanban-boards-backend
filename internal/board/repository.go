package board

import (
	"backend/internal/shared/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"
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
	const op = "board.repository.Get"
	query := "SELECT id, name, description, created_at, updated_at FROM boards WHERE id = $1"
	board := &Board{}
	row := r.storage.QueryRowContext(ctx, query, uuid)
	err := row.Scan(
		&board.ID,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return board, nil
}

func (r *BoardRepository) GetList(
	ctx context.Context,
	filter *BoardGetFilter,
) (*BoardListResult, error) {
	const op = "board.repository.GetList"
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
	const op = "board.repository.Create"
	query := "INSERT INTO boards (name, description, user_id) VALUES ($1, $2, $3)"
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		ErrBoardAlreadyExists,
		board.Name,
		board.Description,
		board.UserID,
	)
}

func (r *BoardRepository) Update(ctx context.Context, board *Board) error {
	const op = "board.repository.Update"
	query := "UPDATE boards SET name = $1, description = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL"
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		ErrBoardNotFound,
		board.Name,
		board.Description,
		board.ID,
	)
}

func (r *BoardRepository) Delete(ctx context.Context, uuid string) error {
	const op = "board.repository.Delete"
	query := "UPDATE boards SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"
	_, err := r.storage.ExecContext(ctx, query, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *BoardRepository) Exists(ctx context.Context, uuid string) (bool, error) {
	const op = "board.repository.Exists"
	var exists bool
	row := r.storage.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM boards WHERE id = $1)",
		uuid,
	)
	err := row.Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *BoardRepository) GetColumnByID(ctx context.Context, column *BoardColumn) (*BoardColumn, error) {
	const op = "board.repository.GetColumnByID"
	query := `
		SELECT id, board_id, position, name, color, created_at
		FROM board_columns bc WHERE id = $1 AND deleted_at IS NULL
		ORDER BY id
	`
	data := &BoardColumn{}
	row := r.storage.QueryRowContext(ctx, query, column.ID)
	err := row.Scan(
		&data.ID,
		&data.BoardID,
		&data.Position,
		&data.Name,
		&data.Color,
		&data.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return data, nil
}

func (r *BoardRepository) GetColumnList(ctx context.Context, uuid string) ([]*BoardColumn, error) {
	const op = "board.repository.GetColumnList"
	columnsRaw := []*BoardColumn{}
	query := `
		SELECT id, board_id, position, name, color, created_at
		FROM board_columns bc WHERE board_id = $1 AND deleted_at IS NULL
		ORDER BY id
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
			&column.Position,
			&column.Name,
			&column.Color,
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
	const op = "board.repository.CreateColumn"
	query := `
		INSERT INTO board_columns (board_id, name, color, position)
		VALUES ($1, $2, $3, $4)
	`
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		ErrBoardNotFound,
		column.BoardID,
		column.Name,
		column.Color,
		column.Position,
	)
}

func (r *BoardRepository) UpdateColumn(ctx context.Context, column *BoardColumn) error {
	const op = "board.repository.UpdateColumn"
	query := `
		UPDATE board_columns
		SET
			name = $1, color = $2, updated_at = NOW()
		WHERE
			board_id = $3 AND id = $4
			and deleted_at is NULL;
	`
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		ErrColumnNotFound,
		column.Name,
		column.Color,
		column.BoardID,
		column.ID,
	)
}

func (r *BoardRepository) DeleteColumn(ctx context.Context, column *BoardColumn) error {
	const op = "board.repository.DeleteColumn"
	tx, err := r.storage.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	defer tx.Rollback()
	query := "UPDATE board_columns SET position = position - 1 WHERE board_id = $1  AND position > $2 AND deleted_at IS NULL"
	_, err = tx.ExecContext(ctx, query, column.BoardID, column.Position)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	query = `
		UPDATE board_columns SET
			position = NULL,
			deleted_at = NOW()
		WHERE board_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	err = utils.OpExec(ctx, tx.ExecContext, op, query, ErrColumnNotFound, column.BoardID, column.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return tx.Commit()
}

func (r *BoardRepository) ExistsColumn(ctx context.Context, uuid string, columnID uint64) (bool, error) {
	const op = "board.repository.ExistsColumn"
	var exists bool
	row := r.storage.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM board_columns WHERE board_id = $1 AND id = $2)",
		uuid,
		columnID,
	)
	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *BoardRepository) GetMaxPositionValue(ctx context.Context, uuid string) (uint64, error) {
	const op = "board.repository.GetMaxPositionValue"
	var maxPosition sql.NullInt64
	var query string
	query = "SELECT MAX(position) FROM board_columns WHERE board_id = $1 AND deleted_at IS NULL"
	row := r.storage.QueryRowContext(
		ctx,
		query,
		uuid,
	)
	err := row.Scan(&maxPosition)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidMaxPositionValue)
	}
	if !maxPosition.Valid {
		return 0, nil
	}
	return uint64(maxPosition.Int64), nil
}

func (r *BoardRepository) MoveColumn(ctx context.Context, id string, columnID, fromPosition, toPosition uint64) error {
	const op = "board.repository.MoveColumn"
	tx, err := r.storage.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	moveQueryToPosition := `UPDATE board_columns SET position = $1 WHERE position = $2 AND board_id = $3 AND deleted_at IS null`
	fromPositionInt := int64(fromPosition)

	err = utils.OpExec(ctx, tx.ExecContext, op, moveQueryToPosition, sql.ErrNoRows, -fromPositionInt, fromPosition, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = utils.OpExec(ctx, tx.ExecContext, op, chooseMoveDirectionQuery(fromPosition, toPosition), sql.ErrNoRows, fromPosition, toPosition, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = utils.OpExec(ctx, tx.ExecContext, op, moveQueryToPosition, sql.ErrNoRows, toPosition, -fromPositionInt, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit()
}

func chooseMoveDirectionQuery(from, to uint64) string {
	var moveQuery string

	if from < to {
		moveQuery = `update board_columns set position = position - 1 where position > $1 AND position <= $2 and board_id = $3 and deleted_at is null` // move to right
	} else {
		moveQuery = `update board_columns set position = position + 1 where position < $1 and position >= $2 and board_id = $3 and deleted_at is null` // move to left
	}
	return moveQuery
}
