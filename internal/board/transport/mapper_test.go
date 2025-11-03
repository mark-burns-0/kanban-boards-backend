package transport

import (
	"backend/internal/board/domain"
	cardDomain "backend/internal/card/domain"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ToBoardGetFilter(t *testing.T) {
	mapper := BoardMapper{}
	name := "test"
	description := "test"

	tests := []struct {
		name     string
		req      *BoardGetFilter
		expected *domain.BoardGetFilter
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "valid data",
			req: &BoardGetFilter{
				UserID:  1,
				PerPage: 1,
				Page:    1,
				FilterFields: &Filters{
					Name:        &name,
					Description: &description,
				},
			},
			expected: &domain.BoardGetFilter{
				UserID:  1,
				PerPage: 1,
				Page:    1,
				FilterFields: &domain.Filters{
					Name:        &name,
					Description: &description,
				},
			},
		},
		{
			name: "only name filter",
			req: &BoardGetFilter{
				UserID:  1,
				PerPage: 1,
				Page:    1,
				FilterFields: &Filters{
					Name: &name,
				},
			},
			expected: &domain.BoardGetFilter{
				UserID:  1,
				PerPage: 1,
				Page:    1,
				FilterFields: &domain.Filters{
					Name: &name,
				},
			},
		},
		{
			name: "only description filter",
			req: &BoardGetFilter{
				UserID:  1,
				PerPage: 1,
				Page:    1,
				FilterFields: &Filters{
					Description: &description,
				},
			},
			expected: &domain.BoardGetFilter{
				UserID:  1,
				PerPage: 1,
				Page:    1,
				FilterFields: &domain.Filters{
					Description: &description,
				},
			},
		},
		{
			name: "without filters",
			req: &BoardGetFilter{
				UserID:       1,
				PerPage:      1,
				Page:         1,
				FilterFields: &Filters{},
			},
			expected: &domain.BoardGetFilter{
				UserID:       1,
				PerPage:      1,
				Page:         1,
				FilterFields: &domain.Filters{},
			},
		},
		{
			name: "pagination only",
			req: &BoardGetFilter{
				UserID:       1,
				PerPage:      1,
				Page:         1,
				FilterFields: nil,
			},
			expected: &domain.BoardGetFilter{
				UserID:       1,
				PerPage:      1,
				Page:         1,
				FilterFields: nil,
			},
		},
		{
			name: "zero values",
			req: &BoardGetFilter{
				UserID:       0,
				PerPage:      0,
				Page:         0,
				FilterFields: &Filters{},
			},
			expected: &domain.BoardGetFilter{
				UserID:       0,
				PerPage:      0,
				Page:         0,
				FilterFields: &domain.Filters{},
			},
		},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("case(%s)", tc.name)
		t.Run(name, func(t *testing.T) {
			actual := mapper.ToBoardGetFilter(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_ToBoard(t *testing.T) {
	mapper := BoardMapper{}
	tests := []struct {
		name     string
		req      *BoardRequest
		expected *domain.Board
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "valid data",
			req: &BoardRequest{
				UserID:      1,
				ID:          "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:        "test",
				Description: "test",
			},
			expected: &domain.Board{
				UserID:      1,
				ID:          "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:        "test",
				Description: "test",
			},
		},
		{
			name: "zero values",
			req: &BoardRequest{
				UserID:      0,
				ID:          "",
				Name:        "",
				Description: "",
			},
			expected: &domain.Board{
				UserID:      0,
				ID:          "",
				Name:        "",
				Description: "",
			},
		},
		{
			name: "empty name",
			req: &BoardRequest{
				UserID:      1,
				ID:          "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:        "",
				Description: "test description",
			},
			expected: &domain.Board{
				UserID:      1,
				ID:          "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:        "",
				Description: "test description",
			},
		},
		{
			name: "empty description",
			req: &BoardRequest{
				UserID:      1,
				ID:          "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:        "test name",
				Description: "",
			},
			expected: &domain.Board{
				UserID:      1,
				ID:          "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:        "test name",
				Description: "",
			},
		},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("case(%s)", tc.name)
		t.Run(name, func(t *testing.T) {
			actual := mapper.ToBoard(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_ToBoardColumn(t *testing.T) {
	mapper := BoardMapper{}
	tests := []struct {
		name     string
		req      *BoardColumnRequest
		expected *domain.BoardColumn
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "valid data",
			req: &BoardColumnRequest{
				ID:      1,
				BoardID: "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:    "test",
				Color:   "color",
			},
			expected: &domain.BoardColumn{
				ID:      1,
				BoardID: "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:    "test",
				Color:   "color",
			},
		},
		{
			name: "zero values",
			req: &BoardColumnRequest{
				ID:      0,
				BoardID: "",
				Name:    "",
				Color:   "",
			},
			expected: &domain.BoardColumn{
				ID:      0,
				BoardID: "",
				Name:    "",
				Color:   "",
			},
		},
		{
			name: "without color",
			req: &BoardColumnRequest{
				ID:      1,
				BoardID: "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:    "test",
			},
			expected: &domain.BoardColumn{
				ID:      1,
				BoardID: "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				Name:    "test",
			},
		},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("case(%s)", tc.name)
		t.Run(name, func(t *testing.T) {
			actual := mapper.ToBoardColumn(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_ToBoardMoveCommand(t *testing.T) {
	mapper := BoardMapper{}
	tests := []struct {
		name     string
		req      *BoardColumnMoveRequest
		expected *domain.BoardMoveCommand
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "valid data",
			req: &BoardColumnMoveRequest{
				BoardID:      "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				ColumnID:     1,
				FromPosition: 1,
				ToPosition:   1,
			},
			expected: &domain.BoardMoveCommand{
				BoardID:      "382a14b1-46f0-4df4-975c-e0d62bd6c358",
				ColumnID:     1,
				FromPosition: 1,
				ToPosition:   1,
			},
		},
		{
			name: "zero values",
			req: &BoardColumnMoveRequest{
				BoardID:      "",
				ColumnID:     0,
				FromPosition: 0,
				ToPosition:   0,
			},
			expected: &domain.BoardMoveCommand{
				BoardID:      "",
				ColumnID:     0,
				FromPosition: 0,
				ToPosition:   0,
			},
		},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("case(%s)", tc.name)
		t.Run(name, func(t *testing.T) {
			actual := mapper.ToBoardMoveCommand(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_ToBoardListResponse(t *testing.T) {
	var zeroValueUint uint64
	mapper := BoardMapper{}
	startedTime := time.Now()
	tests := []struct {
		name     string
		req      *domain.BoardListResult
		expected *BoardListResponse
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "zero values",
			req: &domain.BoardListResult{
				PerPage:     0,
				CurrentPage: 0,
				TotalPages:  0,
				TotalCount:  0,
				NextPage:    &zeroValueUint,
				HasNext:     false,
				HasPrev:     false,
				Data:        []*domain.Board{},
			},
			expected: &BoardListResponse{
				PerPage:     0,
				CurrentPage: 0,
				TotalPages:  0,
				TotalCount:  0,
				NextPage:    &zeroValueUint,
				HasNext:     false,
				HasPrev:     false,
				Data:        []*BoardResponse{},
			},
		},
		{
			name: "valid data",
			req: &domain.BoardListResult{
				PerPage:     1,
				CurrentPage: 1,
				TotalPages:  1,
				TotalCount:  1,
				NextPage:    &zeroValueUint,
				HasNext:     false,
				HasPrev:     true,
				Data: []*domain.Board{
					{
						ID:          "131212313123-123123-12321-3",
						Name:        "Test board",
						Description: "Test description",
						CreatedAt:   startedTime,
						UpdatedAt:   startedTime,
					},
				},
			},
			expected: &BoardListResponse{
				PerPage:     1,
				CurrentPage: 1,
				TotalPages:  1,
				TotalCount:  1,
				NextPage:    &zeroValueUint,
				HasNext:     false,
				HasPrev:     true,
				Data: []*BoardResponse{
					{
						ID:          "131212313123-123123-12321-3",
						Name:        "Test board",
						Description: "Test description",
						CreatedAt:   startedTime,
						UpdatedAt:   startedTime,
					},
				},
			},
		},
		{
			name: "empty data",
			req: &domain.BoardListResult{
				PerPage:     10,
				CurrentPage: 1,
				TotalPages:  0,
				TotalCount:  0,
				NextPage:    nil,
				HasNext:     false,
				HasPrev:     false,
				Data:        []*domain.Board{},
			},
			expected: &BoardListResponse{
				PerPage:     10,
				CurrentPage: 1,
				TotalPages:  0,
				TotalCount:  0,
				NextPage:    nil,
				HasNext:     false,
				HasPrev:     false,
				Data:        []*BoardResponse{},
			},
		},
		{
			name: "nil next page",
			req: &domain.BoardListResult{
				PerPage:     10,
				CurrentPage: 5,
				TotalPages:  5,
				TotalCount:  50,
				NextPage:    nil,
				HasNext:     false,
				HasPrev:     true,
				Data:        []*domain.Board{},
			},
			expected: &BoardListResponse{
				PerPage:     10,
				CurrentPage: 5,
				TotalPages:  5,
				TotalCount:  50,
				NextPage:    nil,
				HasNext:     false,
				HasPrev:     true,
				Data:        []*BoardResponse{},
			},
		},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("case(%s)", tc.name)
		t.Run(name, func(t *testing.T) {
			actual := mapper.ToBoardListResponse(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_ToSingleBoardResponse(t *testing.T) {
	mapper := BoardMapper{}
	testText := "Test text"
	startedTime := time.Now()

	tests := []struct {
		name     string
		req      *domain.BoardWithDetails[cardDomain.CardWithComments]
		expected *SingleBoardResponse[CardWithComments]
	}{
		{
			name:     "nil pointer",
			req:      nil,
			expected: nil,
		},
		{
			name: "zero values",
			req: &domain.BoardWithDetails[cardDomain.CardWithComments]{
				Board: &domain.Board{
					ID:          "",
					Name:        "",
					Description: "",
					CreatedAt:   startedTime,
					UpdatedAt:   startedTime,
				},
				Columns: []*domain.BoardColumn{},
				Cards:   []*cardDomain.CardWithComments{},
			},
			expected: &SingleBoardResponse[CardWithComments]{
				BoardResponse: &BoardResponse{
					ID:          "",
					Name:        "",
					Description: "",
					CreatedAt:   startedTime,
					UpdatedAt:   startedTime,
				},
				Columns: []*BoardColumnResponse{},
				Cards:   []*CardWithComments{},
			},
		},
		{
			name: "valid data",
			req: &domain.BoardWithDetails[cardDomain.CardWithComments]{
				Board: &domain.Board{
					ID:          "93a49b99-a029-4a18-bbbc-c10d91a8c267",
					Name:        "Test name",
					Description: "Test description",
					CreatedAt:   startedTime,
					UpdatedAt:   startedTime,
				},
				Columns: []*domain.BoardColumn{
					{
						ID:        1,
						Position:  1,
						BoardID:   "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Name:      "Test name",
						Color:     "Test color",
						CreatedAt: startedTime,
					},
				},
				Cards: []*cardDomain.CardWithComments{
					{
						ID:          1,
						ColumnID:    1,
						Position:    1,
						BoardID:     "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Text:        testText,
						Description: testText,
						CreatedAt:   startedTime,
						CardProperties: cardDomain.CardProperties{
							Color: testText,
							Tag:   testText,
						},
						Comments: []cardDomain.CardComment{
							{
								ID:        1,
								CardID:    1,
								Text:      testText,
								CreatedAt: startedTime,
							},
						},
					},
				},
			},
			expected: &SingleBoardResponse[CardWithComments]{
				BoardResponse: &BoardResponse{
					ID:          "93a49b99-a029-4a18-bbbc-c10d91a8c267",
					Name:        "Test name",
					Description: "Test description",
					CreatedAt:   startedTime,
					UpdatedAt:   startedTime,
				},
				Columns: []*BoardColumnResponse{
					{
						ID:        1,
						Position:  1,
						BoardID:   "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Name:      "Test name",
						Color:     "Test color",
						CreatedAt: startedTime,
					},
				},
				Cards: []*CardWithComments{
					{
						ID:          1,
						ColumnID:    1,
						Position:    1,
						BoardID:     "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Text:        &testText,
						Description: &testText,
						CreatedAt:   startedTime,
						Properties: &CardProperties{
							Color: &testText,
							Tag:   &testText,
						},
						Comments: []*CardComment{
							{
								ID:        1,
								CardID:    1,
								Text:      testText,
								CreatedAt: startedTime,
							},
						},
					},
				},
			},
		},
		{
			name: "without properties",
			req: &domain.BoardWithDetails[cardDomain.CardWithComments]{
				Board: &domain.Board{
					ID:          "93a49b99-a029-4a18-bbbc-c10d91a8c267",
					Name:        "Test name",
					Description: "Test description",
					CreatedAt:   startedTime,
					UpdatedAt:   startedTime,
				},
				Columns: []*domain.BoardColumn{},
				Cards: []*cardDomain.CardWithComments{
					{
						ID:             1,
						ColumnID:       1,
						Position:       1,
						BoardID:        "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Text:           testText,
						Description:    testText,
						CreatedAt:      startedTime,
						CardProperties: cardDomain.CardProperties{},
						Comments:       []cardDomain.CardComment{},
					},
				},
			},
			expected: &SingleBoardResponse[CardWithComments]{
				BoardResponse: &BoardResponse{
					ID:          "93a49b99-a029-4a18-bbbc-c10d91a8c267",
					Name:        "Test name",
					Description: "Test description",
					CreatedAt:   startedTime,
					UpdatedAt:   startedTime,
				},
				Columns: []*BoardColumnResponse{},
				Cards: []*CardWithComments{
					{
						ID:          1,
						ColumnID:    1,
						Position:    1,
						BoardID:     "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Text:        &testText,
						Description: &testText,
						CreatedAt:   startedTime,
						Properties:  nil,
						Comments:    []*CardComment{},
					},
				},
			},
		},
		{
			name: "unsorted data",
			req: &domain.BoardWithDetails[cardDomain.CardWithComments]{
				Board: &domain.Board{
					ID:          "93a49b99-a029-4a18-bbbc-c10d91a8c267",
					Name:        "Test name",
					Description: "Test description",
					CreatedAt:   startedTime,
					UpdatedAt:   startedTime,
				},
				Columns: []*domain.BoardColumn{
					{
						ID:        2,
						Position:  3,
						BoardID:   "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Name:      "Column 3",
						Color:     "color",
						CreatedAt: startedTime,
					},
					{
						ID:        1,
						Position:  1,
						BoardID:   "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Name:      "Column 1",
						Color:     "color",
						CreatedAt: startedTime,
					},
				},
				Cards: []*cardDomain.CardWithComments{},
			},
			expected: &SingleBoardResponse[CardWithComments]{
				BoardResponse: &BoardResponse{
					ID:          "93a49b99-a029-4a18-bbbc-c10d91a8c267",
					Name:        "Test name",
					Description: "Test description",
					CreatedAt:   startedTime,
					UpdatedAt:   startedTime,
				},
				Columns: []*BoardColumnResponse{
					{
						ID:        1,
						Position:  1,
						BoardID:   "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Name:      "Column 1",
						Color:     "color",
						CreatedAt: startedTime,
					},
					{
						ID:        2,
						Position:  3,
						BoardID:   "93a49b99-a029-4a18-bbbc-c10d91a8c267",
						Name:      "Column 3",
						Color:     "color",
						CreatedAt: startedTime,
					},
				},
				Cards: []*CardWithComments{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapper.ToSingleBoardResponse(tc.req)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
