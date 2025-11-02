package transport

type UserRequest struct {
	Name                 string  `json:"name" validate:"required,min=1,max=255"`
	Email                string  `json:"email" validate:"required,email"`
	Password             *string `json:"password,omitempty" validate:"omitempty,min=8,max=32"`
	PasswordConfirmation *string `json:"password_confirmation,omitempty" validate:"omitempty,min=8,max=32,eqfield=Password"`
}
