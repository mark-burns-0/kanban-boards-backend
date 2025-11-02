package domain

import "time"

type User struct {
	ID                   uint64
	Name                 string
	Email                string
	RefreshToken         string
	Password             string
	PasswordConfirmation string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
}

type RegisterCommand struct {
	Name     string
	Email    string
	Password string
}

type LoginCommand struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
