package transport

import (
	"backend/internal/board/domain"
	cardDomain "backend/internal/card/domain"
	"backend/internal/shared/dto"
	"cmp"
	"slices"
)

type BoardMapper struct{}

func (m *BoardMapper) ToBoardGetFilter(req *BoardGetFilter) *domain.BoardGetFilter {
	if req == nil {
		return nil
	}

	filters := &domain.BoardGetFilter{
		UserID:  req.UserID,
		PerPage: req.PerPage,
		Page:    req.Page,
	}

	if req.FilterFields != nil && req.FilterFields.Name != nil {
		filters.FilterFields.Name = req.FilterFields.Name
	}

	if req.FilterFields != nil && req.FilterFields.Description != nil {
		filters.FilterFields.Description = req.FilterFields.Description
	}

	return filters
}

func (m *BoardMapper) ToBoard(req *BoardRequest) *domain.Board {
	if req == nil {
		return nil
	}

	return &domain.Board{
		ID:          req.ID,
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
	}
}

func (m *BoardMapper) ToBoardColumn(req *BoardColumnRequest) *domain.BoardColumn {
	if req == nil {
		return nil
	}

	return &domain.BoardColumn{
		ID:      req.ID,
		BoardID: req.BoardID,
		Name:    req.Name,
		Color:   req.Color,
	}
}

func (m *BoardMapper) ToBoardMoveCommand(req *BoardColumnMoveRequest) *domain.BoardMoveCommand {
	if req == nil {
		return nil
	}

	return &domain.BoardMoveCommand{
		BoardID:      req.BoardID,
		ColumnID:     req.ColumnID,
		FromPosition: req.FromPosition,
		ToPosition:   req.ToPosition,
	}
}

func (m *BoardMapper) ToBoardListResponse(data *domain.BoardListResult) *BoardListResponse {
	if data == nil {
		return nil
	}

	list := &BoardListResponse{
		PerPage:     data.PerPage,
		CurrentPage: data.CurrentPage,
		NextPage:    data.NextPage,
		TotalPages:  data.TotalPages,
		HasNext:     data.HasNext,
		HasPrev:     data.HasPrev,
		TotalCount:  data.TotalCount,
		Data:        make([]*BoardResponse, 0, len(data.Data)),
	}
	for _, data := range data.Data {
		list.Data = append(list.Data, &BoardResponse{
			ID:          data.ID,
			Name:        data.Name,
			Description: data.Description,
			CreatedAt:   data.CreatedAt,
			UpdatedAt:   data.UpdatedAt,
		})
	}
	return list
}

func (m *BoardMapper) ToSingleBoardResponse(data *domain.BoardWithDetails[cardDomain.CardWithComments]) *SingleBoardResponse[dto.CardWithComments] {
	if data == nil {
		return nil
	}

	board := &SingleBoardResponse[dto.CardWithComments]{
		BoardResponse: &BoardResponse{
			ID:          data.ID,
			Name:        data.Name,
			Description: data.Description,
			CreatedAt:   data.CreatedAt,
			UpdatedAt:   data.UpdatedAt,
		},
		Columns: make([]*BoardColumnResponse, 0, len(data.Columns)),
		Cards:   make([]*dto.CardWithComments, 0, len(data.Cards)),
	}

	for _, column := range data.Columns {
		board.Columns = append(board.Columns, &BoardColumnResponse{
			ID:        column.ID,
			Position:  column.Position,
			BoardID:   column.BoardID,
			Name:      column.Name,
			Color:     column.Color,
			CreatedAt: column.CreatedAt,
		})
	}

	slices.SortStableFunc(board.Columns, func(a, b *BoardColumnResponse) int {
		return cmp.Compare(a.Position, b.Position)
	})

	for _, card := range data.Cards {
		newCard := &dto.CardWithComments{
			ID:          card.ID,
			ColumnID:    card.ColumnID,
			Position:    card.Position,
			BoardID:     card.BoardID,
			Text:        &card.Text,
			Description: &card.Description,
			CreatedAt:   card.CreatedAt,
			Comments:    make([]*dto.CardComment, 0, len(card.Comments)),
		}
		for _, comment := range card.Comments {
			newCard.Comments = append(newCard.Comments, &dto.CardComment{
				ID:        comment.ID,
				CardID:    comment.CardID,
				Text:      comment.Text,
				CreatedAt: comment.CreatedAt,
			})
		}
		if card.CardProperties.Color != "" {
			if newCard.Properties == nil {
				newCard.Properties = &dto.CardProperties{}
			}
			newCard.Properties.Color = &card.CardProperties.Color
		}
		if card.CardProperties.Tag != "" {
			if newCard.Properties == nil {
				newCard.Properties = &dto.CardProperties{}
			}
			newCard.Properties.Tag = &card.CardProperties.Tag
		}
		slices.SortStableFunc(newCard.Comments, func(a, b *dto.CardComment) int {
			return cmp.Compare(a.ID, b.ID)
		})
		board.Cards = append(board.Cards, newCard)
	}
	slices.SortStableFunc(board.Cards, func(a, b *dto.CardWithComments) int {
		return cmp.Compare(a.Position, b.Position)
	})

	return board
}
