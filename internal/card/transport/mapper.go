package transport

import (
	"backend/internal/card/domain"
)

type CardMapper struct{}

func (m *CardMapper) ToCardMoveCommand(req *CardMoveRequest) *domain.CardMoveCommand {
	if req == nil {
		return nil
	}

	return &domain.CardMoveCommand{
		ID:           req.ID,
		FromColumnID: req.FromColumnID,
		ToColumnID:   req.ToColumnID,
		FromPosition: req.FromPosition,
		ToPosition:   req.ToPosition,
		BoardID:      req.BoardID,
	}
}

func (m *CardMapper) ToCard(req *CardRequest) *domain.Card {
	if req == nil {
		return nil
	}
	card := &domain.Card{
		ID:             req.ID,
		ColumnID:       req.ColumnID,
		BoardID:        req.BoardID,
		Text:           req.Text,
		Description:    req.Description,
		CardProperties: domain.CardProperties{},
	}

	card.CardProperties.Color = safeDerefString(req.CardProperties.Color)
	card.CardProperties.Tag = safeDerefString(req.CardProperties.Tag)

	return card
}

func safeDerefString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
