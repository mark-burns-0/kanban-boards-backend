package card

import (
	"context"
	"fmt"
)

type CardGetter interface {
	GetList(context.Context, uint64)
}

type CardCreator interface {
	Create(context.Context, *Card) error
}

type CardUpdater interface {
	Update(context.Context, *Card) error
	MoveToNewPosition(context.Context, string, uint64, uint64, uint64) error
}

type CardDeleter interface {
	Delete(context.Context, *Card) error
}

type CardRepo interface {
	CardCreator
	CardUpdater
	CardDeleter
}

type CardService struct {
	repo CardRepo
}

func NewCardService(repo CardRepo) *CardService {
	return &CardService{
		repo: repo,
	}
}

func (s *CardService) Create(ctx context.Context, request *CardRequest) error {
	op := "card.service.Create"
	card := &Card{
		BoardID:     request.BoardID,
		ColumnID:    request.ColumnID,
		Position:    request.Position,
		Text:        request.Text,
		Description: request.Description,
		cardProperties: cardProperties{
			Color: request.Color,
			Tag:   request.Tag,
		},
	}
	if err := s.repo.Create(ctx, card); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *CardService) Update(ctx context.Context, request *CardRequest) error {
	op := "card.service.Update"
	card := &Card{
		ID:          request.ID,
		ColumnID:    request.ColumnID,
		BoardID:     request.BoardID,
		Position:    request.Position,
		Text:        request.Text,
		Description: request.Description,
		cardProperties: cardProperties{
			Color: request.Color,
			Tag:   request.Tag,
		},
	}
	if err := s.repo.Update(ctx, card); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *CardService) Delete(ctx context.Context, request *CardRequest) error {
	op := "card.service.Delete"
	card := &Card{
		ID:      request.ID,
		BoardID: request.BoardID,
	}
	if err := s.repo.Delete(ctx, card); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *CardService) MoveToNewPosition(ctx context.Context, request *CardMoveRequest) error {
	op := "card.service.MoveToNewPosition"
	if err := s.repo.MoveToNewPosition(
		ctx,
		request.BoardID,
		request.ID,
		request.ToColumnID,
		request.Position,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
