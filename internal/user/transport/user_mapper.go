package transport

import "backend/internal/user/domain"

type UserMapper struct{}

func (m *UserMapper) ToUserDomain(req *UserRequest) *domain.User {
	return &domain.User{
		Name:                 req.Name,
		Email:                req.Email,
		Password:             *req.Password,
		PasswordConfirmation: *req.PasswordConfirmation,
	}
}
