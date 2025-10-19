package user

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

type UserRequest struct {
	Name                 string  `json:"name" validate:"required,min=1,max=255"`
	Email                string  `json:"email" validate:"required,email"`
	Password             *string `json:"password,omitempty" validate:"omitempty,min=8,max=32"`
	PasswordConfirmation *string `json:"password_confirmation,omitempty" validate:"omitempty,min=8,max=32,eqfield=Password"`
}

type UserResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
