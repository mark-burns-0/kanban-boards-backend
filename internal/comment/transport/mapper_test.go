package transport

import (
	"backend/internal/comment/domain"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToComment(t *testing.T) {
	mapper := CommentMapper{}
	tests := []struct {
		name     string
		req      *Comment
		expected *domain.Comment
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "valid data",
			req: &Comment{
				ID:     1,
				CardID: 1,
				UserID: 1,
				Text:   "Test comment",
			},
			expected: &domain.Comment{
				ID:     1,
				CardID: 1,
				UserID: 1,
				Text:   "Test comment",
			},
		},
		{
			name: "empty text",
			req: &Comment{
				ID:     2,
				CardID: 2,
				UserID: 2,
				Text:   "",
			},
			expected: &domain.Comment{
				ID:     2,
				CardID: 2,
				UserID: 2,
				Text:   "",
			},
		},
		{
			name: "zero values",
			req: &Comment{
				ID:     0,
				CardID: 0,
				UserID: 0,
				Text:   "",
			},
			expected: &domain.Comment{
				ID:     0,
				CardID: 0,
				UserID: 0,
				Text:   "",
			},
		},
		{
			name: "long text",
			req: &Comment{
				ID:     3,
				CardID: 3,
				UserID: 3,
				Text:   "This is a very long comment text that might test boundary conditions in the system",
			},
			expected: &domain.Comment{
				ID:     3,
				CardID: 3,
				UserID: 3,
				Text:   "This is a very long comment text that might test boundary conditions in the system",
			},
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case(%s)", test.name)
		t.Run(name, func(t *testing.T) {
			assert.Equal(
				t,
				mapper.ToComment(test.req),
				test.expected,
			)
		})
	}
}
