package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type AuthGetter interface {
	GetByEmail(context.Context, string) (*User, error)
	GetByRefreshToken(context.Context, string) (*User, error)
}

type AuthCreator interface {
	Create(context.Context, *User) error
}

type AuthUpdater interface {
	UpdateRefreshToken(context.Context, uint64, string) error
}

type AuthRepo interface {
	AuthGetter
	AuthCreator
	AuthUpdater
}

type Config interface {
	GetAccessTokenSecret() string
	GetAccessTokenTTL() string
	GetRefreshTokenSecret() string
	GetRefreshTokenTTL() string
	GetBcryptPower() string
}

type AuthService struct {
	authRepo AuthRepo
	config   Config
}

func NewAuthService(authRepo AuthRepo, config Config) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		config:   config,
	}
}

func (r *AuthService) Register(ctx context.Context, userRequest *UserCreateRequest) error {
	op := "auth.service.Register"
	user, _ := r.authRepo.GetByEmail(ctx, userRequest.Email)
	if user != nil {
		return fmt.Errorf("user with email %s already exists", userRequest.Email)
	}

	power, err := strconv.Atoi(r.config.GetBcryptPower())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), power)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newUser := User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: string(hashedPassword),
	}

	return r.authRepo.Create(ctx, &newUser)
}

func (r *AuthService) Login(ctx context.Context, userRequest *UserLoginRequest) (*TokensResponse, error) {
	op := "auth.service.Login"
	user, _ := r.authRepo.GetByEmail(ctx, userRequest.Email)
	if user == nil {
		return nil, fmt.Errorf("%s: user with email %s not found", op, userRequest.Email)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		return nil, fmt.Errorf("%s: invalid password", op)
	}

	accessToken, err := r.generateToken(user, AccessTokenType)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	refreshToken, err := r.generateToken(user, RefreshTokenType)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = r.authRepo.UpdateRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (r *AuthService) RefreshToken(ctx context.Context, token string) (*TokensResponse, error) {
	op := "auth.service.RefreshToken"
	user, err := r.authRepo.GetByRefreshToken(ctx, token)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if user == nil {
		return nil, fmt.Errorf("%s: user with refresh token %s not found", op, token)
	}

	accessToken, err := r.generateToken(user, AccessTokenType)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	refreshToken, err := r.generateToken(user, RefreshTokenType)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.authRepo.UpdateRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (r *AuthService) generateToken(user *User, tokenType string) (string, error) {
	op := "auth.service.generateToken"
	if tokenType == "" {
		return "", fmt.Errorf("token type is required")
	}

	tokenSecret, ttl, err := r.getTokenConfig(tokenType)
	if err != nil {
		return "", fmt.Errorf("%s: failed to get token config: %w", op, err)
	}

	key := []byte(tokenSecret)
	duration, err := time.ParseDuration(ttl)
	if err != nil {
		return "", fmt.Errorf("%s: invalid TTL duration for %s token: %w", op, tokenType, err)
	}

	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": map[string]any{
				"id": user.ID,
			},
			"exp": time.Now().Add(duration).Unix(),
		})

	signedString, err := claims.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign %s token: %w", op, tokenType, err)
	}

	return signedString, nil
}

func (r *AuthService) getTokenConfig(tokenType string) (string, string, error) {
	op := "auth.service.getTokenConfig"
	switch tokenType {
	case AccessTokenType:
		return r.config.GetAccessTokenSecret(), r.config.GetAccessTokenTTL(), nil
	case RefreshTokenType:
		return r.config.GetRefreshTokenSecret(), r.config.GetRefreshTokenTTL(), nil
	default:
		return "", "", fmt.Errorf("%s: unknown token type: %s", op, tokenType)
	}
}
