package domain

import (
	"backend/internal/shared/domain/events"
	"backend/internal/shared/errors"
	"context"
)

type CardBoardEventHandler struct {
	repo CardRepo
}

func NewCardBoardEventHandler(repo CardRepo) *CardBoardEventHandler {
	return &CardBoardEventHandler{
		repo: repo,
	}
}

func (h *CardBoardEventHandler) Handle(ctx context.Context, event events.Event) error {
	switch e := event.(type) {
	case events.BoardDeletedEvent:
		exists, err := h.repo.CardExistsInBoard(ctx, e.BoardID)
		if err != nil {
			return err
		}
		if exists {
			return errors.ErrBoardHasCards
		}
	case events.ColumnDeletedEvent:
		exists, err := h.repo.CardExistsInColumn(ctx, e.ColumnID)
		if err != nil {
			return err
		}
		if exists {
			return errors.ErrColumnHasCards
		}
	}
	return nil
}
