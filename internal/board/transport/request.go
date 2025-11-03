package transport

import "time"

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

type CardProperties struct {
	Color *string `json:"color,omitempty"`
	Tag   *string `json:"tag,omitempty"`
}

type CardWithComments struct {
	ID          uint64          `json:"id"`
	ColumnID    uint64          `json:"column_id"`
	Position    uint64          `json:"position"`
	BoardID     string          `json:"board_id"`
	Text        *string         `json:"text"`
	Description *string         `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	Properties  *CardProperties `json:"properties,omitempty"`
	Comments    []*CardComment  `json:"comments"`
}

type CardComment struct {
	ID        uint64    `json:"id"`
	CardID    uint64    `json:"card_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
