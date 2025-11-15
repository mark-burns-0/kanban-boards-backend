package repository

import (
	"backend/internal/board/domain"
	"backend/internal/shared/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	existsBoardQuery  = "SELECT EXISTS(SELECT 1 FROM boards WHERE id = $1 AND deleted_at IS NULL)"
	existsColumnQuery = "SELECT EXISTS(SELECT 1 FROM board_columns WHERE board_id = $1 AND id = $2 AND deleted_at IS NULL)"
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
	GetDB() *sql.DB
	Close() error
}

type BoardRepository struct {
	storage Storage
}

func NewBoardRepository(storage Storage) *BoardRepository {
	return &BoardRepository{
		storage: storage,
	}
}

func (r *BoardRepository) Get(ctx context.Context, info *domain.Board) (*domain.Board, error) {
	const op = "board.repository.Get"
	query := "SELECT id, name, description, created_at, updated_at FROM boards WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL"
	board := &domain.Board{}

	row := r.storage.QueryRowContext(ctx, query, info.ID, info.UserID)
	err := row.Scan(
		&board.ID,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBoardNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return board, nil
}

func (r *BoardRepository) GetList(
	ctx context.Context,
	filter *domain.BoardGetFilter,
) (*domain.BoardListResult, error) {
	const op = "board.repository.GetList"
	boards := []*domain.Board{}
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
		board := &domain.Board{}
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
	result := &domain.BoardListResult{
		Data:       boards,
		TotalCount: count,
	}
	return result, nil
}

func buildQuery(filter *domain.BoardGetFilter) (string, []any) {
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
	filter *domain.Filters,
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

func (r *BoardRepository) Create(ctx context.Context, board *domain.Board) error {
	const op = "board.repository.Create"
	query := "INSERT INTO boards (name, description, user_id) VALUES ($1, $2, $3)"
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		domain.ErrBoardAlreadyExists,
		board.Name,
		board.Description,
		board.UserID,
	)
}

func (r *BoardRepository) Update(ctx context.Context, board *domain.Board) error {
	const op = "board.repository.Update"
	query := "UPDATE boards SET name = $1, description = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL"
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		domain.ErrBoardNotFound,
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
	var err error

	if exists, err = utils.ExistsQueryWrapper(
		ctx,
		r.storage,
		existsBoardQuery,
		uuid,
	); err != nil {
		return exists, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}
