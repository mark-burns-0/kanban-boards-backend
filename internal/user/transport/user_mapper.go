package transport

import "backend/internal/user/domain"

type UserMapper struct{}

func (m *UserMapper) ToUserDomain(req *UserRequest) *domain.User {
	if req == nil {
		return nil
	}

	return &domain.User{
		Name:                 req.Name,
		Email:                req.Email,
		Password:             safeDerefString(req.Password),
		PasswordConfirmation: safeDerefString(req.PasswordConfirmation),
	}
}

func safeDerefString(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}
