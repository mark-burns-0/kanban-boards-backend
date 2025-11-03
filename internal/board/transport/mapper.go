package transport

import (
	"backend/internal/board/domain"
	cardDomain "backend/internal/card/domain"
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

	if req.FilterFields != nil {
		if filters.FilterFields == nil {
			filters.FilterFields = &domain.Filters{}
		}

		if req.FilterFields.Name != nil {
			filters.FilterFields.Name = req.FilterFields.Name
		}

		if req.FilterFields.Description != nil {
			filters.FilterFields.Description = req.FilterFields.Description
		}
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

func (m *BoardMapper) ToSingleBoardResponse(data *domain.BoardWithDetails[cardDomain.CardWithComments]) *SingleBoardResponse[CardWithComments] {
	if data == nil {
		return nil
	}

	return &SingleBoardResponse[CardWithComments]{
		BoardResponse: m.toBoardResponse(data),
		Columns:       m.mapAndSortColumns(data.Columns),
		Cards:         m.mapAndSortCards(data.Cards),
	}
}

func (m *BoardMapper) toBoardResponse(data *domain.BoardWithDetails[cardDomain.CardWithComments]) *BoardResponse {
	return &BoardResponse{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func (m *BoardMapper) mapAndSortColumns(columns []*domain.BoardColumn) []*BoardColumnResponse {
	mapped := make([]*BoardColumnResponse, 0, len(columns))
	for _, column := range columns {
		mapped = append(mapped, &BoardColumnResponse{
			ID:        column.ID,
			Position:  column.Position,
			BoardID:   column.BoardID,
			Name:      column.Name,
			Color:     column.Color,
			CreatedAt: column.CreatedAt,
		})
	}

	slices.SortFunc(mapped, func(a, b *BoardColumnResponse) int {
		return cmp.Compare(a.Position, b.Position)
	})

	return mapped
}

func (m *BoardMapper) mapAndSortCards(cards []*cardDomain.CardWithComments) []*CardWithComments {
	mapped := make([]*CardWithComments, 0, len(cards))
	for _, card := range cards {
		mapped = append(mapped, m.mapCardWithComments(card))
	}

	slices.SortFunc(mapped, func(a, b *CardWithComments) int {
		return cmp.Compare(a.ID, b.ID)
	})

	return mapped
}

func (m *BoardMapper) mapCardWithComments(card *cardDomain.CardWithComments) *CardWithComments {
	mapped := &CardWithComments{
		ID:          card.ID,
		ColumnID:    card.ColumnID,
		Position:    card.Position,
		BoardID:     card.BoardID,
		Text:        &card.Text,
		Description: &card.Description,
		CreatedAt:   card.CreatedAt,
		Comments:    m.mapAndSortComments(card.Comments),
	}

	m.mapCardProperties(card, mapped)
	return mapped
}

func (m *BoardMapper) mapCardProperties(card *cardDomain.CardWithComments, mapped *CardWithComments) {
	if card.CardProperties.Color == "" && card.CardProperties.Tag == "" {
		return
	}

	mapped.Properties = &CardProperties{}
	if card.CardProperties.Color != "" {
		mapped.Properties.Color = &card.CardProperties.Color
	}
	if card.CardProperties.Tag != "" {
		mapped.Properties.Tag = &card.CardProperties.Tag
	}
}

func (m *BoardMapper) mapAndSortComments(comments []cardDomain.CardComment) []*CardComment {
	mapped := make([]*CardComment, 0, len(comments))
	for _, comment := range comments {
		mapped = append(mapped, &CardComment{
			ID:        comment.ID,
			CardID:    comment.CardID,
			Text:      comment.Text,
			CreatedAt: comment.CreatedAt,
		})
	}

	slices.SortFunc(mapped, func(a, b *CardComment) int {
		return cmp.Compare(a.ID, b.ID)
	})

	return mapped
}
