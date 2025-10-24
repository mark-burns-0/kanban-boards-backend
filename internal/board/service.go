package board

import (
	"context"
	"fmt"
)

type BoardCreator interface {
	Create(ctx context.Context, board *Board) error
}

type BoardUpdater interface {
	Update(ctx context.Context, board *Board) error
	MoveToColumn(ctx context.Context, id string, columnID, fromPosition, toPosition uint64) error
}

type BoardDeleter interface {
	Delete(ctx context.Context, uuid string) error
}

type BoardGetter interface {
	Get(ctx context.Context, uuid string) (*Board, error)
	GetList(ctx context.Context, userID uint64) ([]*Board, error)
}

type BoardRepo interface {
	BoardCreator
	BoardUpdater
	BoardDeleter
	BoardGetter
}

type BoardService struct {
	repo BoardRepo
}

func NewBoardService(repo BoardRepo) *BoardService {
	return &BoardService{
		repo: repo,
	}
}

func (s *BoardService) Create(ctx context.Context, board *BoardRequest) error {
	op := "board.service.Create"
	data := &Board{
		UserID:      board.UserID,
		Name:        board.Name,
		Description: board.Description,
	}
	if err := s.repo.Create(ctx, data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *BoardService) Update(ctx context.Context, board *BoardRequest) error {
	op := "board.service.Update"
	data := &Board{
		ID:          board.ID,
		UserID:      board.UserID,
		Name:        board.Name,
		Description: board.Description,
	}
	if err := s.repo.Update(ctx, data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *BoardService) Delete(ctx context.Context, boardUUID string) error {
	op := "board.service.Delete"
	if err := s.repo.Delete(ctx, boardUUID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
