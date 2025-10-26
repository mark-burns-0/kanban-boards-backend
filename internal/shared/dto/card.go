package dto

import "time"

type CardProperties struct {
	Color *string `json:"color,omitempty"`
	Tag   *string `json:"tag,omitempty"`
}

type CardWithComments struct {
	ID          *uint64         `json:"id"`
	ColumnID    *uint64         `json:"column_id"`
	Position    *uint64         `json:"position"`
	BoardID     *string         `json:"board_id"`
	Text        *string         `json:"text"`
	Description *string         `json:"description"`
	CreatedAt   *time.Time      `json:"created_at"`
	Properties  *CardProperties `json:"properties,omitempty"`
	Comments    []*CardComment  `json:"comments"`
}

type CardComment struct {
	ID        *uint64    `json:"id"`
	CardID    *uint64    `json:"card_id"`
	Text      *string    `json:"text"`
	CreatedAt *time.Time `json:"created_at"`
}
