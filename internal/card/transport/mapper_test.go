package transport

import (
	"backend/internal/card/domain"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToCardMoveCommand(t *testing.T) {
	mapper := CardMapper{}
	tests := []struct {
		name     string
		req      *CardMoveRequest
		expected *domain.CardMoveCommand
	}{
		{
			name: "zero values",
			req: &CardMoveRequest{
				ID:           0,
				FromColumnID: 0,
				ToColumnID:   0,
				FromPosition: 0,
				ToPosition:   0,
				BoardID:      "",
			},
			expected: &domain.CardMoveCommand{
				ID:           0,
				FromColumnID: 0,
				ToColumnID:   0,
				FromPosition: 0,
				ToPosition:   0,
				BoardID:      "",
			},
		},
		{
			name: "empty board ID",
			req: &CardMoveRequest{
				ID:           1,
				FromColumnID: 2,
				ToColumnID:   3,
				FromPosition: 4,
				ToPosition:   5,
				BoardID:      "",
			},
			expected: &domain.CardMoveCommand{
				ID:           1,
				FromColumnID: 2,
				ToColumnID:   3,
				FromPosition: 4,
				ToPosition:   5,
				BoardID:      "",
			},
		},
		{
			name: "valid data",
			req: &CardMoveRequest{
				ID:           1,
				FromColumnID: 2,
				ToColumnID:   3,
				FromPosition: 4,
				ToPosition:   5,
				BoardID:      "e102c99e-651c-44e1-bff1-c4a22e3134ce",
			},
			expected: &domain.CardMoveCommand{
				ID:           1,
				FromColumnID: 2,
				ToColumnID:   3,
				FromPosition: 4,
				ToPosition:   5,
				BoardID:      "e102c99e-651c-44e1-bff1-c4a22e3134ce",
			},
		},
		{
			name:     "nil pointer check",
			req:      nil,
			expected: nil,
		},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("case (%s)", tc.name)
		t.Run(name, func(t *testing.T) {
			actual := mapper.ToCardMoveCommand(tc.req)
			assert.Equal(t, actual, tc.expected)
		})
	}
}
func Test_ToCard(t *testing.T) {
	mapper := CardMapper{}
	var color, tag *string
	color = new(string)
	tag = new(string)
	*color = "color"
	*tag = "tag"

	tests := []struct {
		name     string
		req      *CardRequest
		expected *domain.Card
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "main information",
			req: &CardRequest{
				ID:          1,
				ColumnID:    1,
				BoardID:     "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:        "test text",
				Description: "test description",
			},
			expected: &domain.Card{
				ID:          1,
				ColumnID:    1,
				BoardID:     "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:        "test text",
				Description: "test description",
			},
		},
		{
			name: "with color and tag properties",
			req: &CardRequest{
				ID:          1,
				ColumnID:    1,
				BoardID:     "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:        "test text",
				Description: "test description",
				CardProperties: CardProperties{
					Color: color,
					Tag:   tag,
				},
			},
			expected: &domain.Card{
				ID:          1,
				ColumnID:    1,
				BoardID:     "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:        "test text",
				Description: "test description",
				CardProperties: domain.CardProperties{
					Color: "color",
					Tag:   "tag",
				},
			},
		},
		{
			name: "with color property only",
			req: &CardRequest{
				ID:          1,
				ColumnID:    1,
				BoardID:     "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:        "test text",
				Description: "test description",
				CardProperties: CardProperties{
					Color: color,
				},
			},
			expected: &domain.Card{
				ID:          1,
				ColumnID:    1,
				BoardID:     "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:        "test text",
				Description: "test description",
				CardProperties: domain.CardProperties{
					Color: "color",
				},
			},
		},
		{
			name: "with tag property only",
			req: &CardRequest{
				ID:          1,
				ColumnID:    1,
				BoardID:     "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:        "test text",
				Description: "test description",
				CardProperties: CardProperties{
					Tag: tag,
				},
			},
			expected: &domain.Card{
				ID:          1,
				ColumnID:    1,
				BoardID:     "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:        "test text",
				Description: "test description",
				CardProperties: domain.CardProperties{
					Tag: "tag",
				},
			},
		},
		{
			name: "with nil properties pointer",
			req: &CardRequest{
				ID:             1,
				ColumnID:       1,
				BoardID:        "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:           "test text",
				Description:    "test description",
				CardProperties: CardProperties{},
			},
			expected: &domain.Card{
				ID:             1,
				ColumnID:       1,
				BoardID:        "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:           "test text",
				Description:    "test description",
				CardProperties: domain.CardProperties{},
			},
		},
		{
			name: "with nil property values",
			req: &CardRequest{
				ID:          1,
				ColumnID:    1,
				BoardID:     "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:        "test text",
				Description: "test description",
				CardProperties: CardProperties{
					Color: nil,
					Tag:   nil,
				},
			},
			expected: &domain.Card{
				ID:             1,
				ColumnID:       1,
				BoardID:        "e102c99e-651c-44e1-bff1-c4a22e3134ce",
				Text:           "test text",
				Description:    "test description",
				CardProperties: domain.CardProperties{},
			},
		},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("case (%s)", tc.name)
		t.Run(name, func(t *testing.T) {
			actual := mapper.ToCard(tc.req)
			assert.Equal(t, actual, tc.expected)
		})
	}
}
