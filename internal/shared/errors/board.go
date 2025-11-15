package errors

import "errors"

var (
	ErrBoardHasCards  = errors.New("board has cards")
	ErrColumnHasCards = errors.New("column has cards")
)
