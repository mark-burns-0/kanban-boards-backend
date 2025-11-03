package domain

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

func (r *AuthService) Register(ctx context.Context, req *RegisterCommand) error {
	const op = "auth.service.Register"
	user, _ := r.authRepo.GetByEmail(ctx, req.Email)
	if user != nil {
		return fmt.Errorf("%s: %w", op, UserAlreadyExistsErr{
			email: req.Email,
		})
	}

	power, err := strconv.Atoi(r.config.GetBcryptPower())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), power)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newUser := User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return r.authRepo.Create(ctx, &newUser)
}

func (r *AuthService) Login(ctx context.Context, req *LoginCommand) (*Tokens, error) {
	const op = "auth.service.Login"
	user, _ := r.authRepo.GetByEmail(ctx, req.Email)
	if user == nil {
		err := EmailNotFoundErr{
			email: req.Email,
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, InvalidPasswordErr{})
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

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (r *AuthService) RefreshToken(ctx context.Context, token string) (*Tokens, error) {
	const op = "auth.service.RefreshToken"
	user, err := r.authRepo.GetByRefreshToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if user == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
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

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (r *AuthService) generateToken(user *User, tokenType string) (string, error) {
	const op = "auth.service.generateToken"
	if tokenType == "" {
		return "", tokenTypeRequiredErr{}
	}

	tokenSecret, ttl, err := r.getTokenConfig(tokenType)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrFailedToGetTokenConfig)
	}

	key := []byte(tokenSecret)
	duration, err := time.ParseDuration(ttl)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, invalidDurationTTLErr{
			tokenType, err,
		})
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
	const op = "auth.service.getTokenConfig"
	switch tokenType {
	case AccessTokenType:
		return r.config.GetAccessTokenSecret(), r.config.GetAccessTokenTTL(), nil
	case RefreshTokenType:
		return r.config.GetRefreshTokenSecret(), r.config.GetRefreshTokenTTL(), nil
	default:
		return "", "", fmt.Errorf("%s: %w", op, unkownTokenTypeErr{tokenType})
	}
}
