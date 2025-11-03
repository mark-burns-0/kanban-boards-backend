package transport

import (
	"backend/internal/user/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToRegisterCommand(t *testing.T) {
	mapper := AuthMapper{}

	tests := []struct {
		name     string
		req      *UserRegisterRequest
		expected *domain.RegisterCommand
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "valid data",
			req: &UserRegisterRequest{
				Name:                 "test",
				Email:                "test@gmail.com",
				Password:             "testtest",
				PasswordConfirmation: "testtest",
			},
			expected: &domain.RegisterCommand{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "testtest",
			},
		},
		{
			name: "empty strings",
			req: &UserRegisterRequest{
				Name:                 "",
				Email:                "",
				Password:             "",
				PasswordConfirmation: "",
			},
			expected: &domain.RegisterCommand{
				Name:     "",
				Email:    "",
				Password: "",
			},
		},
		{
			name: "without password confirmation",
			req: &UserRegisterRequest{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "testtest",
			},
			expected: &domain.RegisterCommand{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "testtest",
			},
		},
		{
			name: "with different passwords", // если это допустимо в маппере
			req: &UserRegisterRequest{
				Name:                 "test",
				Email:                "test@gmail.com",
				Password:             "password1",
				PasswordConfirmation: "password2",
			},
			expected: &domain.RegisterCommand{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "password1",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapper.ToRegisterCommand(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_ToLoginCommand(t *testing.T) {
	mapper := AuthMapper{}

	tests := []struct {
		name     string
		req      *UserLoginRequest
		expected *domain.LoginCommand
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "valid data",
			req: &UserLoginRequest{
				Email:    "test@gmail.com",
				Password: "testtest",
			},
			expected: &domain.LoginCommand{
				Email:    "test@gmail.com",
				Password: "testtest",
			},
		},
		{
			name: "empty credentials",
			req: &UserLoginRequest{
				Email:    "",
				Password: "",
			},
			expected: &domain.LoginCommand{
				Email:    "",
				Password: "",
			},
		},
		{
			name: "only email",
			req: &UserLoginRequest{
				Email: "test@gmail.com",
			},
			expected: &domain.LoginCommand{
				Email: "test@gmail.com",
			},
		},
		{
			name: "only password",
			req: &UserLoginRequest{
				Password: "testtest",
			},
			expected: &domain.LoginCommand{
				Password: "testtest",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapper.ToLoginCommand(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_ToResponseTokens(t *testing.T) {
	mapper := AuthMapper{}

	const (
		accessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjIyMTU0NzIsInN1YiI6eyJpZCI6Mn19.Si5_OZwOu4mP715pcpN5hehUcfgLDmPoYg02p8HgVpk"
		refreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ3Nzg2NzIsInN1YiI6eyJpZCI6Mn19.BQOAkW0kctj1xQzqUzMTTb2-CxOF_oi-QXzsdfkP730"
	)

	tests := []struct {
		name     string
		req      *domain.Tokens
		expected *TokensResponse
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "valid data",
			req: &domain.Tokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
			expected: &TokensResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		},
		{
			name: "only access token",
			req: &domain.Tokens{
				AccessToken: accessToken,
			},
			expected: &TokensResponse{
				AccessToken: accessToken,
			},
		},
		{
			name: "only refresh token",
			req: &domain.Tokens{
				RefreshToken: refreshToken,
			},
			expected: &TokensResponse{
				RefreshToken: refreshToken,
			},
		},
		{
			name: "empty tokens",
			req: &domain.Tokens{
				AccessToken:  "",
				RefreshToken: "",
			},
			expected: &TokensResponse{
				AccessToken:  "",
				RefreshToken: "",
			},
		},
		{
			name: "empty access token",
			req: &domain.Tokens{
				AccessToken:  "",
				RefreshToken: refreshToken,
			},
			expected: &TokensResponse{
				AccessToken:  "",
				RefreshToken: refreshToken,
			},
		},
		{
			name: "empty refresh token",
			req: &domain.Tokens{
				AccessToken:  accessToken,
				RefreshToken: "",
			},
			expected: &TokensResponse{
				AccessToken:  accessToken,
				RefreshToken: "",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapper.ToResponseTokens(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
