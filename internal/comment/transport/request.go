package transport

type Comment struct {
	ID     uint64 `json:"id,omitempty" validate:"number"`
	CardID uint64 `json:"card_id,omitempty" validate:"required,number"`
	UserID uint64 `json:"user_id,omitempty" validate:"required,number"`
	Text   string `json:"text" validate:"required,min=1,max=4096"`
}
