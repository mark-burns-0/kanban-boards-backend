package transport

import "backend/internal/comment/domain"

type CommentMapper struct{}

func (m *CommentMapper) ToComment(comment *Comment) *domain.Comment {
	return &domain.Comment{
		ID:     comment.ID,
		UserID: comment.UserID,
		CardID: comment.CardID,
		Text:   comment.Text,
	}
}
