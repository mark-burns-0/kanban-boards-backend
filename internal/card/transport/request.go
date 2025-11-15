package transport

type CardRequest struct {
	ID             uint64
	ColumnID       uint64 `json:"column_id" validate:"required,min=1"`
	BoardID        string `json:"board_id" validate:"required,uuid"`
	Text           string `json:"text" validate:"required,min=1,max=255"`
	Description    string `json:"description" validate:"required,min=1,max=255"`
	CardProperties `json:"properties"`
}

type CardMoveRequest struct {
	ID           uint64 `json:"id" validate:"required,min=1"`
	FromColumnID uint64 `json:"from_column_id" validate:"required,min=1"`
	ToColumnID   uint64 `json:"to_column_id" validate:"required,min=1"`
	FromPosition uint64 `json:"from_position" validate:"required,min=1"`
	ToPosition   uint64 `json:"to_position" validate:"required,min=1"`
	BoardID      string `json:"board_id" validate:"required,uuid"`
}

type CardProperties struct {
	Color *string `json:"color,omitempty" validate:"omitnil,hexcolor,max=255"`
	Tag   *string `json:"tag,omitempty" validate:"omitnil,max=255"`
}
