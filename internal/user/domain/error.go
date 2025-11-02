package domain

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrFailedToGetTokenConfig = errors.New("failed to get token config")
	ErrInvalidPassword        = errors.New("invalid password")
)

type UserAlreadyExistsErr struct {
	email string
}

func (e UserAlreadyExistsErr) Error() string {
	return fmt.Sprintf("user with email %s already exists", e.email)
}

func (e UserAlreadyExistsErr) Unwrap() error {
	return ErrUserAlreadyExists
}

type EmailNotFoundErr struct {
	email string
}

func (e EmailNotFoundErr) Error() string {
	return fmt.Sprintf("user with email %s not found", e.email)
}

func (e EmailNotFoundErr) Unwrap() error {
	return ErrUserNotFound
}

type InvalidPasswordErr struct {
}

func (e InvalidPasswordErr) Error() string {
	return "invalid password"
}

func (e InvalidPasswordErr) Unwrap() error {
	return ErrInvalidPassword
}

type tokenTypeRequiredErr struct{}

func (e tokenTypeRequiredErr) Error() string {
	return "token type is required"
}

func (e tokenTypeRequiredErr) Unwrap() error {
	return e
}

type unkownTokenTypeErr struct {
	tokenType string
}

func (e unkownTokenTypeErr) Error() string {
	return fmt.Sprintf("unknown token type: %s", e.tokenType)
}

func (e unkownTokenTypeErr) Unwrap() error {
	return e
}

type invalidDurationTTLErr struct {
	tokenType string
	err       error
}

func (e invalidDurationTTLErr) Error() string {
	return fmt.Sprintf("invalid duration TTL for token type %s", e.tokenType)
}

func (e invalidDurationTTLErr) Unwrap() error {
	return e.err
}
