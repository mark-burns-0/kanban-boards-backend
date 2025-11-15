package domain

import "time"

type BoardColumn struct {
	ID        uint64
	Position  uint64
	BoardID   string
	Name      string
	Color     string
	CreatedAt time.Time
}
