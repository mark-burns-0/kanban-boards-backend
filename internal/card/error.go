package card

import (
	"errors"
)

var (
	ErrCardAlreadyExists = errors.New("card already exists")
	ErrCardNotFound      = errors.New("card not found")
	ErrColumnNotExist    = errors.New("column not exists")
)
