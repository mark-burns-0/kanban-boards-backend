package repository

import (
	"backend/internal/board/domain"
	"backend/internal/shared/utils"
	"context"
	"database/sql"
	"fmt"
)

func (r *BoardRepository) GetColumnByID(ctx context.Context, column *domain.BoardColumn) (*domain.BoardColumn, error) {
	const op = "board.repository.GetColumnByID"
	query := `
		SELECT id, board_id, position, name, color, created_at
		FROM board_columns bc WHERE id = $1 AND deleted_at IS NULL
		ORDER BY id
	`
	data := &domain.BoardColumn{}
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

func (r *BoardRepository) GetColumnList(ctx context.Context, uuid string) ([]*domain.BoardColumn, error) {
	const op = "board.repository.GetColumnList"
	columnsRaw := []*domain.BoardColumn{}
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
		column := &domain.BoardColumn{}
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

func (r *BoardRepository) CreateColumn(ctx context.Context, column *domain.BoardColumn) error {
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
		domain.ErrBoardNotFound,
		column.BoardID,
		column.Name,
		column.Color,
		column.Position,
	)
}

func (r *BoardRepository) UpdateColumn(ctx context.Context, column *domain.BoardColumn) error {
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
		domain.ErrColumnNotFound,
		column.Name,
		column.Color,
		column.BoardID,
		column.ID,
	)
}

func (r *BoardRepository) DeleteColumn(ctx context.Context, column *domain.BoardColumn) error {
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
	err = utils.OpExec(ctx, tx.ExecContext, op, query, domain.ErrColumnNotFound, column.BoardID, column.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return tx.Commit()
}

func (r *BoardRepository) ExistsColumn(ctx context.Context, uuid string, columnID uint64) (bool, error) {
	const op = "board.repository.ExistsColumn"
	var exists bool
	var err error

	if exists, err = utils.ExistsQueryWrapper(
		ctx,
		r.storage,
		existsColumnQuery,
		uuid,
		columnID,
	); err != nil {
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
		return 0, fmt.Errorf("%s: %w", op, domain.ErrInvalidMaxPositionValue)
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
