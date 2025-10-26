package card

import (
	"time"
)

type Card struct {
	ID          uint64
	ColumnID    uint64
	Position    uint64
	BoardID     string
	Text        string
	Description string
	cardProperties
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CardRequest struct {
	ID             uint64
	ColumnID       uint64 `json:"column_id" validate:"required,min=1"`
	Position       uint64 `json:"position" validate:"required,min=1"`
	BoardID        string `json:"board_id" validate:"required,uuid"`
	Text           string `json:"text" validate:"required,min=1,max=255"`
	Description    string `json:"description" validate:"required,min=1,max=255"`
	cardProperties `json:"properties"`
}

type CardMoveRequest struct {
	ID           uint64 `json:"id" validate:"required,min=1"`
	FromColumnID uint64 `json:"from_column_id" validate:"required,min=1"`
	ToColumnID   uint64 `json:"to_column_id" validate:"required,min=1"`
	Position     uint64 `json:"position" validate:"required,min=1"`
	BoardID      string `json:"board_id" validate:"required,uuid"`
}

type cardProperties struct {
	Color string `json:"color,omitempty" validate:"hexcolor,max=255"`
	Tag   string `json:"tag,omitempty" validate:"max=255"`
}
