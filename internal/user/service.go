package user

import "context"

type UserFinder interface {
	GetByID(context.Context, uint64) (*User, error)
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

func (s *UserService) Current(ctx context.Context) (*User, error) {
	return nil, nil
}

func (s *UserService) Update(ctx context.Context, userRequest *UserRequest) error {
	return nil
}
