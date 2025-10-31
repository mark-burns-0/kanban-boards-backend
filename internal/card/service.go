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
	GetMaxColumnPosition(ctx context.Context, boardUUID string, columnID uint64) (uint64, error)
	GetById(ctx context.Context, card *Card) (*Card, error)
}

type CardCreator interface {
	Create(context.Context, *Card) error
	Exists(ctx context.Context, card *Card) (bool, error)
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

type BoardRepo interface {
	ExistsColumn(ctx context.Context, uuid string, columnID uint64) (bool, error)
}

type CardService struct {
	repo      CardRepo
	boardRepo BoardRepo
}

func NewCardService(repo CardRepo, boardRepo BoardRepo) *CardService {
	return &CardService{
		repo:      repo,
		boardRepo: boardRepo,
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

func (s *CardService) Create(ctx context.Context, req *CardRequest) error {
	const op = "card.service.Create"
	maxPosition, err := s.repo.GetMaxColumnPosition(ctx, req.BoardID, req.ColumnID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	card := &Card{
		BoardID:     req.BoardID,
		ColumnID:    req.ColumnID,
		Text:        req.Text,
		Position:    maxPosition + 1,
		Description: req.Description,
		cardProperties: cardProperties{
			Color: req.Color,
			Tag:   req.Tag,
		},
	}
	if err := s.repo.Create(ctx, card); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *CardService) Update(ctx context.Context, req *CardRequest) error {
	const op = "card.service.Update"
	card := &Card{
		ID:          req.ID,
		ColumnID:    req.ColumnID,
		BoardID:     req.BoardID,
		Text:        req.Text,
		Description: req.Description,
		cardProperties: cardProperties{
			Color: req.Color,
			Tag:   req.Tag,
		},
	}
	exists, err := s.repo.Exists(ctx, card)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrCardNotFound)
	}
	if err := s.repo.Update(ctx, card); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *CardService) Delete(ctx context.Context, req *CardRequest) error {
	const op = "card.service.Delete"
	card := &Card{
		ID:      req.ID,
		BoardID: req.BoardID,
	}
	exists, err := s.repo.Exists(ctx, card)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrCardNotFound)
	}
	card, err = s.repo.GetById(ctx, card)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := s.repo.Delete(ctx, card); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *CardService) MoveToNewPosition(ctx context.Context, req *CardMoveRequest) error {
	const op = "card.service.MoveToNewPosition"
	exists, err := s.repo.Exists(ctx, &Card{
		ID:      req.ID,
		BoardID: req.BoardID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrCardNotFound)
	}
	isFromToColunExists, _ := s.boardRepo.ExistsColumn(ctx, req.BoardID, req.FromColumnID)
	isToColumnExists, _ := s.boardRepo.ExistsColumn(ctx, req.BoardID, req.ToColumnID)
	if !isFromToColunExists || !isToColumnExists {
		return fmt.Errorf("%s: %w", op, ErrColumnNotExist)
	}
	if req.FromColumnID == req.ToColumnID && req.FromPosition == req.ToPosition {
		return nil
	}
	maxValue, err := s.repo.GetMaxColumnPosition(ctx, req.BoardID, req.ToColumnID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if req.ToPosition < 1 || req.ToPosition > maxValue {
		req.ToPosition = 1
	}
	if err := s.repo.MoveToNewPosition(
		ctx,
		req.BoardID,
		req.ID,
		req.FromColumnID,
		req.ToColumnID,
		req.FromPosition,
		req.ToPosition,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
