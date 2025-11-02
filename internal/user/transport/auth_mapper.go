package transport

import "backend/internal/user/domain"

type AuthMapper struct{}

func (m *AuthMapper) ToRegisterCommand(req *UserRegisterRequest) *domain.RegisterCommand {
	return &domain.RegisterCommand{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func (m *AuthMapper) ToLoginCommand(req *UserLoginRequest) *domain.LoginCommand {
	return &domain.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (m *AuthMapper) ToResponseTokens(tokens *domain.Tokens) *TokensResponse {
	return &TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
