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

type cardProperties struct {
	Color string `json:"color,omitempty" validate:"max=255"`
	Tag   string `json:"tag,omitempty" validate:"max=255"`
}
