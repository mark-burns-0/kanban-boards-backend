package transport

import (
	"backend/internal/comment/domain"
)

type CommentMapper struct{}

func (m *CommentMapper) ToComment(req *Comment) *domain.Comment {
	if req == nil {
		return nil
	}

	return &domain.Comment{
		ID:     req.ID,
		UserID: req.UserID,
		CardID: req.CardID,
		Text:   req.Text,
	}
}
