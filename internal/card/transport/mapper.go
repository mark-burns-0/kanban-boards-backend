package transport

import "backend/internal/card/domain"

type CardMapper struct{}

func (m *CardMapper) ToCardMoveCommand(req *CardMoveRequest) *domain.CardMoveCommand {
	return &domain.CardMoveCommand{
		ID:           req.ID,
		FromColumnID: req.FromColumnID,
		ToColumnID:   req.ToColumnID,
		FromPosition: req.FromPosition,
		ToPosition:   req.ToPosition,
		BoardID:      req.BoardID,
	}
}

func (m *CardMapper) CardRequestToCard(req *CardRequest) *domain.Card {
	card := &domain.Card{
		ID:          req.ID,
		ColumnID:    req.ColumnID,
		BoardID:     req.BoardID,
		Text:        req.Text,
		Description: req.Description,
	}

	if req.cardProperties != nil && req.cardProperties.Color != nil {
		card.Color = *req.cardProperties.Color
	}

	if req.cardProperties != nil && req.cardProperties.Tag != nil {
		card.Tag = *req.cardProperties.Tag
	}

	return card
}
