package utils

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func OpExec(
	ctx context.Context, storage Storage, op string, query string, mainErr error, params ...any,
) error {
	result, err := storage.ExecContext(
		ctx,
		query,
		params...,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("%s: %w", op, mainErr)
	}
	return nil
}
