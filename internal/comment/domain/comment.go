package domain

import "time"

type Comment struct {
	ID        uint64
	CardID    uint64
	UserID    uint64
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
