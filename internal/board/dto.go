package board

import (
	"time"
)

type Board struct {
	ID          uint64
	UserID      uint64
	Name        string
	Description string
	IsPublic    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type BoardCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public,omitempty"`
}
