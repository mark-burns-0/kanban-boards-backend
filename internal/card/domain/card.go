package domain

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
	CardProperties
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CardWithComments struct {
	ID          uint64
	ColumnID    uint64
	Position    uint64
	BoardID     string
	Text        string
	Description string
	CardProperties
	Comments  []CardComment
	CreatedAt time.Time
}

type CardComment struct {
	ID        uint64
	CardID    uint64
	Text      string
	CreatedAt time.Time
}

type CardMoveCommand struct {
	ID           uint64
	FromColumnID uint64
	ToColumnID   uint64
	FromPosition uint64
	ToPosition   uint64
	BoardID      string
}

type CardProperties struct {
	Color string
	Tag   string
}

func (cp *CardProperties) Scan(value interface{}) error {
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

func (cp CardProperties) Value() (driver.Value, error) {
	if cp.Color == "" && cp.Tag == "" {
		return "{}", nil
	}
	return json.Marshal(cp)
}
