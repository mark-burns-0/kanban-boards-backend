package domain

import (
	"backend/internal/shared/domain/events"
	"backend/internal/shared/errors"
	"context"
)

type CommentCardEventHandler struct {
	repo CommentRepo
}

func NewCommentCardEventHandler(repo CommentRepo) *CommentCardEventHandler {
	return &CommentCardEventHandler{
		repo: repo,
	}
}

func (h *CommentCardEventHandler) Handle(ctx context.Context, event events.Event) error {
	switch e := event.(type) {
	case events.CardDeletedEvent:
		exists, err := h.repo.CommentExistsInCard(ctx, e.CardID)
		if err != nil {
			return err
		}
		if exists {
			return errors.ErrCardHasComments
		}
	}
	return nil
}
