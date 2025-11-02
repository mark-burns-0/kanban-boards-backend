package board

import "errors"

var (
	ErrBoardAlreadyExists      = errors.New("board already exists")
	ErrBoardNotFound           = errors.New("board not found")
	ErrColumnNotFound          = errors.New("column not found")
	ErrInvalidPosition         = errors.New("invalid position")
	ErrInvalidMaxPositionValue = errors.New("invalid max position value")
)
