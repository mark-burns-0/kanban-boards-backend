package utils

import (
	"context"
	"database/sql"
	"fmt"
)

type exexSQLFunc func(ctx context.Context, query string, args ...any) (sql.Result, error)

func OpExec(
	ctx context.Context, exec exexSQLFunc, op string, query string, mainErr error, params ...any,
) error {
	result, err := exec(
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

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, mainErr)
	}
	return nil
}

// ExistsQueryWrapper выполняет EXISTS SQL запрос и возвращает результат.
// Пробрасывает ошибки без изменения для последующего wrapping в доменном коде.
func ExistsQueryWrapper(
	ctx context.Context,
	db interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	},
	query string,
	args ...any,
) (bool, error) {
	row := db.QueryRowContext(
		ctx,
		query,
		args...,
	)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
