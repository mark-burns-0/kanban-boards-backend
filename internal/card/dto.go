package card

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type cardProperties map[string]string

func (cp cardProperties) Value() (driver.Value, error) {
	if cp == nil {
		return nil, nil
	}
	return json.Marshal(cp)
}

func (cp *cardProperties) Scan(value any) error {
	if value == nil {
		*cp = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, cp)
}

type Card struct {
	ID          uint64
	BoardID     uint64
	ColumnID    uint64
	Text        string
	Description string
	Position    uint64
	AssignedTo  uint64
	Properties  cardProperties
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type CardRequest struct {
	BoardID     uint64 `json:"board_id" validate:"required,min=1"`
	ColumnID    uint64 `json:"column_id" validate:"required,min=1"`
	Text        string `json:"text" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=1,max=255"`
	Position    uint64 `json:"position" validate:"required,min=1"`
}
