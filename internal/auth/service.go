package auth

import "context"

type AuthGetter interface {
	GetByEmail(context.Context, string) (*User, error)
	GetByRefreshToken(context.Context, string) (*User, error)
}

type AuthCreator interface {
	Create(context.Context, *UserCreateRequest) error
}

type AuthRepo interface {
	AuthGetter
	AuthCreator
}

type AuthService struct {
	authRepo AuthRepo
}

func NewAuthService(authRepo AuthRepo) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

func (r *AuthService) Register(userRequest *UserCreateRequest) error {
	return nil
}

func (r *AuthService) Login(userRequest *UserLoginRequest) (*UserResponse, error) {
	return nil, nil
}

func (r *AuthService) RefreshToken(ctx context.Context) (*TokensResponse, error) {
	return nil, nil
}
