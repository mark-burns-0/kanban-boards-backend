package board

import (
	"context"
	"fmt"
	"math"
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
	GetList(ctx context.Context, filter *BoardGetFilter) (*BoardListResult, error)
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

func (s *BoardService) GetList(
	ctx context.Context, filter *BoardGetFilter,
) (*BoardListResponse, error) {
	op := "board.service.GetList"
	rawResp, err := s.repo.GetList(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	boards := make([]*BoardResponse, 0, len(rawResp.Data))
	for _, board := range rawResp.Data {
		boards = append(
			boards,
			&BoardResponse{
				ID:          board.ID,
				Name:        board.Name,
				Description: board.Description,
				CreatedAt:   board.CreatedAt,
				UpdatedAt:   board.UpdatedAt,
			},
		)
	}

	totalPages := uint64(math.Ceil(float64(rawResp.TotalCount) / float64(filter.PerPage)))
	hasPrev := filter.Page > 1
	hasNext := filter.Page < totalPages
	var nextPage *uint64
	if hasNext {
		page := filter.Page + 1
		nextPage = &page
	}
	response := &BoardListResponse{
		Data:        boards,
		PerPage:     filter.PerPage,
		CurrentPage: filter.Page,
		NextPage:    nextPage,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
		TotalCount:  rawResp.TotalCount,
	}
	return response, nil
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
