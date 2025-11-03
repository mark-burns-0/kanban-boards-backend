package transport

import "time"

type BoardResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SingleBoardResponse[T any] struct {
	*BoardResponse
	Columns []*BoardColumnResponse `json:"columns"`
	Cards   []*T                   `json:"cards"`
}

type BoardColumnResponse struct {
	ID        uint64    `json:"id"`
	Position  uint64    `json:"position"`
	BoardID   string    `json:"board_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}
type BoardListResponse struct {
	Data        []*BoardResponse `json:"data"`
	PerPage     uint64           `json:"per_page"`
	CurrentPage uint64           `json:"current_page"`
	TotalCount  uint64           `json:"total_count"`
	TotalPages  uint64           `json:"total_pages"`
	NextPage    *uint64          `json:"next_page"`
	HasNext     bool             `json:"has_next"`
	HasPrev     bool             `json:"has_prev"`
}
