package comment

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

type CommentRequest struct {
	ID     uint64 `json:"id,omitempty" validate:"number"`
	CardID uint64 `json:"card_id,omitempty" validate:"required,number"`
	UserID uint64 `json:"user_id,omitempty" validate:"required,number"`
	Text   string `json:"text" validate:"required,min=1,max=4096"`
}
