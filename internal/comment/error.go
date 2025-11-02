package comment

import "errors"

var (
	ErrCardNotFound         = errors.New("card not found")
	ErrCommentAlreadyExists = errors.New("comment already exists")
)
