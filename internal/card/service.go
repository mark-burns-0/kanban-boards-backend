package card

import (
	"backend/internal/shared/dto"
	"cmp"
	"context"
	"fmt"
	"slices"
)

type CardGetter interface {
	GetListWithComments(ctx context.Context, boardID string) ([]*CardWithComments, error)
}

type CardCreator interface {
	Create(context.Context, *Card) error
}

type CardUpdater interface {
	Update(context.Context, *Card) error
	MoveToNewPosition(ctx context.Context, boardID string, cardID, fromColumnID, toColumnID, cardFromPosition, cardToPosition uint64) error
}

type CardDeleter interface {
	Delete(context.Context, *Card) error
}

type CardRepo interface {
	CardGetter
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

func (s *CardService) GetListWithComments(ctx context.Context, boardID string) ([]*dto.CardWithComments, error) {
	const op = "card.service.GetListWithComments"
	raws, err := s.repo.GetListWithComments(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	response := make([]*dto.CardWithComments, 0, len(raws))

	for _, raw := range raws {
		card := &dto.CardWithComments{
			ID:          raw.ID,
			ColumnID:    raw.ColumnID,
			Position:    raw.Position,
			BoardID:     raw.BoardID,
			Text:        raw.Text,
			Description: raw.Description,
			CreatedAt:   raw.CreatedAt,
			Properties: &dto.CardProperties{
				Color: raw.cardProperties.Color,
				Tag:   raw.cardProperties.Tag,
			},
			Comments: make([]*dto.CardComment, 0, len(raw.Comments)),
		}
		var comment *dto.CardComment
		for _, rawComment := range raw.Comments {
			comment = &dto.CardComment{
				ID:        rawComment.ID,
				CardID:    rawComment.CardID,
				Text:      rawComment.Text,
				CreatedAt: rawComment.CreatedAt,
			}
			card.Comments = append(card.Comments, comment)
		}

		response = append(response, card)
	}
	slices.SortFunc(response, func(a, b *dto.CardWithComments) int {
		return cmp.Compare(*a.ID, *b.ID)
	})
	return response, nil
}

func (s *CardService) Create(ctx context.Context, request *CardRequest) error {
	const op = "card.service.Create"
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
	const op = "card.service.Update"
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
	const op = "card.service.Delete"
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
	const op = "card.service.MoveToNewPosition"
	if err := s.repo.MoveToNewPosition(
		ctx,
		request.BoardID,
		request.ID,
		request.FromColumnID,
		request.ToColumnID,
		request.FromPosition,
		request.ToPosition,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
