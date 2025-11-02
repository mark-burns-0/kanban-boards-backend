package domain

import "time"

type Board struct {
	ID          string
	Name        string
	Description string
	UserID      uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type BoardColumn struct {
	ID        uint64
	Position  uint64
	BoardID   string
	Name      string
	Color     string
	CreatedAt time.Time
}

type BoardListResult struct {
	Data       []*Board
	TotalCount uint64
}

type BoardWithDetails[T any] struct {
	*Board
	Columns []*BoardColumn `json:"columns"`
	Cards   []*T           `json:"cards"`
}

type BoardMoveCommand struct {
	BoardID      string
	ColumnID     uint64
	FromPosition uint64
	ToPosition   uint64
}

type BoardGetFilter struct {
	UserID       uint64
	PerPage      uint64
	Page         uint64
	FilterFields *Filters
}

type Filters struct {
	Name        *string
	Description *string
}
