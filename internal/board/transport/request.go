package transport

type BoardRequest struct {
	UserID      uint64
	ID          string
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"required,min=2,max=1000"`
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

type BoardColumnRequest struct {
	ID      uint64
	BoardID string `json:"board_id" validate:"required,uuid"`
	Name    string `json:"name" validate:"required,min=2"`
	Color   string `json:"color" validate:"required,hexcolor"`
}

type BoardColumnMoveRequest struct {
	BoardID      string `json:"board_id" validate:"required,uuid"`
	ColumnID     uint64 `json:"column_id" validate:"required,min=1"`
	FromPosition uint64 `json:"from_position" validate:"required,min=1"`
	ToPosition   uint64 `json:"to_position" validate:"required,min=1"`
}
