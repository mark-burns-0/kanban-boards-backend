package auth

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthGetter interface {
	GetByEmail(context.Context, string) (*User, error)
	GetByRefreshToken(context.Context, string) (*User, error)
}

type AuthCreator interface {
	Create(context.Context, *User) error
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
	user, _ := r.authRepo.GetByEmail(context.TODO(), userRequest.Email)
	if user != nil {
		return fmt.Errorf("user with email %s already exists", userRequest.Email)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 12)
	if err != nil {
		return err
	}
	newUser := User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: string(hashedPassword),
	}
	return r.authRepo.Create(context.TODO(), &newUser)
}

func (r *AuthService) Login(userRequest *UserLoginRequest) (*UserResponse, error) {
	user, _ := r.authRepo.GetByEmail(context.TODO(), userRequest.Email)
	if user == nil {
		return nil, fmt.Errorf("user with email %s not found", userRequest.Email)
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}
	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (r *AuthService) RefreshToken(ctx context.Context) (*TokensResponse, error) {
	return nil, nil
}
