package board

import (
	"backend/internal/card/domain"
	"cmp"
	"context"
	"fmt"
	"math"
	"slices"

	"golang.org/x/sync/errgroup"
)

type BoardCreator interface {
	Create(ctx context.Context, board *Board) error
	CreateColumn(ctx context.Context, column *BoardColumn) error
}

type BoardUpdater interface {
	Update(ctx context.Context, board *Board) error
	UpdateColumn(ctx context.Context, column *BoardColumn) error
	MoveColumn(ctx context.Context, id string, columnID, fromPosition, toPosition uint64) error
}

type BoardDeleter interface {
	Delete(ctx context.Context, uuid string) error
	DeleteColumn(ctx context.Context, column *BoardColumn) error
}

type BoardGetter interface {
	Get(ctx context.Context, uuid string) (*Board, error)
	GetList(ctx context.Context, filter *BoardGetFilter) (*BoardListResult, error)
	GetColumnByID(ctx context.Context, column *BoardColumn) (*BoardColumn, error)
	GetColumnList(ctx context.Context, uuid string) ([]*BoardColumn, error)
	GetMaxPositionValue(ctx context.Context, uuid string) (uint64, error)
	Exists(ctx context.Context, uuid string) (bool, error)
	ExistsColumn(ctx context.Context, uuid string, columnID uint64) (bool, error)
}

type BoardRepo interface {
	BoardCreator
	BoardUpdater
	BoardDeleter
	BoardGetter
}

type CardService interface {
	GetListWithComments(ctx context.Context, boardID string) ([]*domain.CardWithComments, error)
	Create(ctx context.Context, req *domain.Card) error
	Update(ctx context.Context, req *domain.Card) error
	Delete(ctx context.Context, req *domain.Card) error
	MoveToNewPosition(ctx context.Context, req *domain.CardMoveCommand) error
}

type BoardService struct {
	repo        BoardRepo
	cardService CardService
}

func NewBoardService(repo BoardRepo, cardService CardService) *BoardService {
	return &BoardService{
		repo:        repo,
		cardService: cardService,
	}
}

func (s *BoardService) GetList(
	ctx context.Context, filter *BoardGetFilter,
) (*BoardListResponse, error) {
	const op = "board.service.GetList"
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

func (s *BoardService) GetByUUID(ctx context.Context, boardUUID string) (*SingleBoardResponse[domain.CardWithComments], error) {
	const op = "board.service.GetByUUID"
	var (
		rawBoard   *Board
		rawColumns []*BoardColumn
		cards      []*domain.CardWithComments
	)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		result, err := s.cardService.GetListWithComments(ctx, boardUUID)
		cards = result
		return err
	})
	eg.Go(func() error {
		result, err := s.repo.Get(ctx, boardUUID)
		rawBoard = result
		return err
	})
	eg.Go(func() error {
		result, err := s.repo.GetColumnList(ctx, boardUUID)
		rawColumns = result
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var columns []*BoardColumnResponse
	for _, rawColumn := range rawColumns {
		column := &BoardColumnResponse{
			ID:        rawColumn.ID,
			Name:      rawColumn.Name,
			Color:     rawColumn.Color,
			Position:  rawColumn.Position,
			BoardID:   rawColumn.BoardID,
			CreatedAt: rawColumn.CreatedAt,
		}
		columns = append(columns, column)
	}
	slices.SortFunc(columns, func(a, b *BoardColumnResponse) int {
		return cmp.Compare(a.ID, b.ID)
	})
	board := &SingleBoardResponse[domain.CardWithComments]{
		BoardResponse: &BoardResponse{
			ID:          rawBoard.ID,
			Name:        rawBoard.Name,
			Description: rawBoard.Description,
			CreatedAt:   rawBoard.CreatedAt,
			UpdatedAt:   rawBoard.UpdatedAt,
		},
		Cards:   cards,
		Columns: columns,
	}
	return board, nil
}

func (s *BoardService) Create(ctx context.Context, board *BoardRequest) error {
	const op = "board.service.Create"
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
	const op = "board.service.Update"
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
	const op = "board.service.Delete"
	exists, err := s.repo.Exists(ctx, boardUUID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrBoardNotFound)
	}
	if err := s.repo.Delete(ctx, boardUUID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *BoardService) CreateColumn(ctx context.Context, req *BoardColumnRequest) error {
	const op = "board.service.CreateColumn"
	maxVal, err := s.repo.GetMaxPositionValue(ctx, req.BoardID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	column := &BoardColumn{
		BoardID:  req.BoardID,
		Name:     req.Name,
		Color:    req.Color,
		Position: maxVal + 1,
	}
	if err := s.repo.CreateColumn(ctx, column); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *BoardService) UpdateColumn(ctx context.Context, req *BoardColumnRequest) error {
	const op = "board.service.UpdateColumn"
	column := &BoardColumn{
		ID:      req.ID,
		BoardID: req.BoardID,
		Name:    req.Name,
		Color:   req.Color,
	}
	if err := s.repo.UpdateColumn(ctx, column); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *BoardService) DeleteColumn(ctx context.Context, req *BoardColumnRequest) error {
	const op = "board.service.DeleteColumn"
	column := &BoardColumn{
		ID:      req.ID,
		BoardID: req.BoardID,
	}
	column, err := s.repo.GetColumnByID(ctx, column)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := s.repo.DeleteColumn(ctx, column); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *BoardService) MoveColumn(ctx context.Context, req *BoardColumnMoveRequest) error {
	const op = "board.service.MoveColumn"
	exists, err := s.repo.ExistsColumn(ctx, req.BoardID, req.ColumnID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrColumnNotFound)
	}
	if req.FromPosition == req.ToPosition {
		return nil
	}
	maxValue, err := s.repo.GetMaxPositionValue(ctx, req.BoardID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if req.FromPosition > maxValue || req.ToPosition > maxValue+1 {
		return fmt.Errorf("%s: %w", op, ErrInvalidPosition)
	}
	if err := s.repo.MoveColumn(
		ctx,
		req.BoardID,
		req.ColumnID,
		req.FromPosition,
		req.ToPosition,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
