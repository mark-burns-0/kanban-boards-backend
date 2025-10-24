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
