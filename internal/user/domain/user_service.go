package domain

import (
	"backend/internal/user/transport"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordMismatch = fmt.Errorf("password mismatch")
)

type UserFinder interface {
	Get(context.Context, uint64) (*User, error)
}

type UserUpdater interface {
	Update(context.Context, *User) error
}

type UserRepo interface {
	UserFinder
	UserUpdater
}

type Config interface {
	GetAccessTokenSecret() string
	GetAccessTokenTTL() string
	GetRefreshTokenSecret() string
	GetRefreshTokenTTL() string
	GetBcryptPower() string
}

type UserService struct {
	userRepo UserRepo
	config   Config
}

func NewUserService(userRepo UserRepo, config Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *UserService) Current(ctx context.Context, userID uint64) (*transport.UserResponse, error) {
	const op = "user.service.Current"
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &transport.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: &user.CreatedAt,
	}, nil
}

func (s *UserService) Update(ctx context.Context, userRequest *transport.UserRequest, userID uint64) error {
	const op = "user.service.Update"
	user := &User{
		ID:    userID,
		Name:  userRequest.Name,
		Email: userRequest.Email,
	}
	if userRequest.Password != nil {
		if userRequest.PasswordConfirmation == nil || *userRequest.Password != *userRequest.PasswordConfirmation {
			return ErrPasswordMismatch
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(*userRequest.Password), 12)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		user.Password = string(hashed)
	}
	return s.userRepo.Update(ctx, user)
}
