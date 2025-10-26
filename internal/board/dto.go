package board

import (
	"time"
)

type Board struct {
	ID          string
	UserID      uint64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type BoardRequest struct {
	UserID      uint64
	ID          string
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"required,min=2,max=1000"`
}

type BoardColumRequest struct {
	BoardID  string `json:"board_id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Position uint64 `json:"position"`
}

type BoardResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SingleBoardResponse struct {
	*BoardResponse
	Columns []BoardColumnResponse `json:"columns"`
	Cards   []any                 `json:"cards"`
}

type BoardGetFilter struct {
	UserID       uint64   `validate:"required,min=1"`
	PerPage      uint64   `json:"per_page" validate:"required,min=1,max=200"`
	Page         uint64   `json:"page" validate:"required,min=1"`
	FilterFields *Filters `json:"filters,omitempty"`
}

type Filters struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}

type BoardListResponse struct {
	Data        []*BoardResponse `json:"data"`
	PerPage     uint64           `json:"per_page"`
	CurrentPage uint64           `json:"current_page"`
	NextPage    *uint64          `json:"next_page"`
	TotalPages  uint64           `json:"total_pages"`
	HasNext     bool             `json:"has_next"`
	HasPrev     bool             `json:"has_prev"`
	TotalCount  uint64           `json:"total_count"`
}

type BoardListResult struct {
	Data       []*Board
	TotalCount uint64
}

type BoardColumnResponse struct {
	BoardID  string `json:"board_id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Position uint64 `json:"position"`
}

type BoardColumn struct {
	ID        uint64
	Position  uint64
	BoardID   string
	Name      string
	Color     string
	CreatedAt time.Time
}
