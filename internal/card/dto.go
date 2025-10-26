package card

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

type CardWithComments struct {
	ID          *uint64
	ColumnID    *uint64
	Position    *uint64
	BoardID     *string
	Text        *string
	Description *string
	*cardProperties
	CreatedAt *time.Time
	Comments  []*CardComment
}

type CardComment struct {
	ID        *uint64
	CardID    *uint64
	Text      *string
	CreatedAt *time.Time
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
	Color *string `json:"color,omitempty" validate:"omitnil,hexcolor,max=255"`
	Tag   *string `json:"tag,omitempty" validate:"omitnil,max=255"`
}

func (cp *cardProperties) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("CardProperties.Scan: expected []byte, got %T", value)
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, cp)
}

func (cp cardProperties) Value() (driver.Value, error) {
	if cp.Color == nil && cp.Tag == nil {
		return "{}", nil
	}
	return json.Marshal(cp)
}
