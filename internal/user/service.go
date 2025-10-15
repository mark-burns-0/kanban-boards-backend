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

type UserService struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Current(ctx context.Context) (*User, error) {
	return nil, nil
}

func (s *UserService) Update(ctx context.Context, userRequest *UserRequest) error {
	return nil
}
