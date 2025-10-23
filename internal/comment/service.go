package comment

import (
	"context"
	"fmt"
)

type CommentCreator interface {
	Create(ctx context.Context, comment *Comment) error
}

type CommentUpdater interface {
	Update(ctx context.Context, comment *Comment) error
}

type CommentDeleter interface {
	Delete(ctx context.Context, commentID uint64) error
}

type CommentRepo interface {
	CommentCreator
	CommentUpdater
	CommentDeleter
}

type CommentService struct {
	repository CommentRepo
}

func NewCommentService(repository CommentRepo) *CommentService {
	return &CommentService{
		repository: repository,
	}
}

func (s *CommentService) Create(ctx context.Context, req *CommentRequest) error {
	op := "comment.service.Create"

	comment := &Comment{
		UserID: req.UserID,
		CardID: req.CardID,
		Text:   req.Text,
	}
	err := s.repository.Create(ctx, comment)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *CommentService) Update(ctx context.Context, req *CommentRequest) error {
	op := "comment.service.Update"

	comment := &Comment{
		ID:     req.ID,
		CardID: req.CardID,
		Text:   req.Text,
	}

	if err := s.repository.Update(ctx, comment); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *CommentService) Delete(ctx context.Context, commentID uint64) error {
	op := "comment.service.Delete"

	if err := s.repository.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
