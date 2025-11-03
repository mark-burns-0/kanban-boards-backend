package transport

import (
	"backend/internal/user/domain"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToUserDomain(t *testing.T) {
	mapper := UserMapper{}
	var password, passwordConfirmation *string
	password = new(string)
	passwordConfirmation = new(string)
	*password = "test"
	*passwordConfirmation = "test"

	tests := []struct {
		name     string
		req      *UserRequest
		expected *domain.User
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "valid data",
			req: &UserRequest{
				Name:                 "Test",
				Email:                "test@gmail.com",
				Password:             password,
				PasswordConfirmation: passwordConfirmation,
			},
			expected: &domain.User{
				Name:                 "Test",
				Email:                "test@gmail.com",
				Password:             "test",
				PasswordConfirmation: "test",
			},
		},
		{
			name: "without password",
			req: &UserRequest{
				Name:  "Test",
				Email: "test@gmail.com",
			},
			expected: &domain.User{
				Name:  "Test",
				Email: "test@gmail.com",
			},
		},
		{
			name: "empty strings",
			req: &UserRequest{
				Name:  "",
				Email: "",
			},
			expected: &domain.User{
				Name:  "",
				Email: "",
			},
		},
		{
			name: "with password but without confirmation",
			req: &UserRequest{
				Name:     "Test",
				Email:    "test@gmail.com",
				Password: password,
			},
			expected: &domain.User{
				Name:     "Test",
				Email:    "test@gmail.com",
				Password: "test",
			},
		},
		{
			name: "with confirmation but without password",
			req: &UserRequest{
				Name:                 "Test",
				Email:                "test@gmail.com",
				PasswordConfirmation: passwordConfirmation,
			},
			expected: &domain.User{
				Name:                 "Test",
				Email:                "test@gmail.com",
				PasswordConfirmation: "test",
			},
		},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("case(%s)", tc.name)
		t.Run(name, func(t *testing.T) {
			actual := mapper.ToUserDomain(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
