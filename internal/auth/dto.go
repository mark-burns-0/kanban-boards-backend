package auth

import "time"

type User struct {
	ID           uint64
	Name         string
	Email        string
	Password     string
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type UserCreateRequest struct {
	Name                 string `json:"name" validate:"required,min=1,max=255"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,max=32"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,max=32,eqfield=Password"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
