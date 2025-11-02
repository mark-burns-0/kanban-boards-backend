package transport

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=1,max=255"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,max=32"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,max=32,eqfield=Password"`
}
