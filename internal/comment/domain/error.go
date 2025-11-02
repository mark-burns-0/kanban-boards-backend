package domain

import "errors"

var (
	ErrCardNotFound         = errors.New("card not found")
	ErrCommentNotFound      = errors.New("comment not found")
	ErrCommentAlreadyExists = errors.New("comment already exists")
)
