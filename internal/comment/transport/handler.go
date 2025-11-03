package transport

import (
	"backend/internal/comment/domain"
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey    = "userID"
	CardIDKey    = "card_id"
	CommentIDKey = "comment_id"
)

const (
	CreatedMessage = "created"
	UpdatedMessage = "updated"
)

type CommentService interface {
	Create(ctx context.Context, req *domain.Comment) error
	Update(ctx context.Context, req *domain.Comment) error
	Delete(ctx context.Context, commentID uint64) error
}

type CommentHandler struct {
	validaotr     http.Validator
	lang          http.LangMessage
	service       CommentService
	commentMapper *CommentMapper
}

func NewCommentHandler(
	validator http.Validator,
	lang http.LangMessage,
	service CommentService,
) *CommentHandler {
	return &CommentHandler{
		validaotr:     validator,
		lang:          lang,
		service:       service,
		commentMapper: &CommentMapper{},
	}
}

func (h *CommentHandler) Create(c *fiber.Ctx) error {
	const op = "comment.transport.handler.Create"
	comment, err := utils.ParseBody[Comment](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	cardID, err := strconv.ParseUint(c.Params(CardIDKey), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": domain.ErrCardNotFound})
	}
	comment.CardID = cardID

	if validationErrors, statusCode, err := h.validaotr.ValidateStruct(c, comment); validationErrors != nil {
		if err != nil {
			slog.Error("validator error",
				slog.String("op", op),
				slog.Any("err", err),
			)
			return c.Status(statusCode).JSON(fiber.Map{"error": "Validation error"})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}

	if err := h.service.Create(c.Context(), h.commentMapper.ToComment(comment)); err != nil {
		switch {
		case errors.Is(err, domain.ErrCommentAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Comment already exists"})
		case errors.Is(err, domain.ErrCardNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Card not found"})
		}
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("error", err),
		)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), CreatedMessage),
		},
	)
}

func (h *CommentHandler) Update(c *fiber.Ctx) error {
	const op = "comment.transport.handler.Update"
	comment, err := utils.ParseBody[Comment](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	commentID := c.Params(CommentIDKey)
	commentIDUint64, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": domain.ErrCommentNotFound})
	}

	cardID, err := strconv.ParseUint(c.Params(CardIDKey), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": domain.ErrCardNotFound})
	}
	comment.CardID = cardID
	comment.ID = commentIDUint64

	if validationErrors, statusCode, err := h.validaotr.ValidateStruct(c, comment); validationErrors != nil {
		if err != nil {
			slog.Error("validator error",
				slog.String("op", op),
				slog.Any("err", err),
			)
			return c.Status(statusCode).JSON(fiber.Map{"error": "Validation error"})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}

	if err := h.service.Update(c.Context(), h.commentMapper.ToComment(comment)); err != nil {
		switch {
		case errors.Is(err, domain.ErrCardNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Card not found"})
		}
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("error", err),
		)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Server error"})
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), UpdatedMessage),
		},
	)
}

func (h *CommentHandler) Delete(c *fiber.Ctx) error {
	commentID := c.Params(CommentIDKey)
	commentIDUint64, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": domain.ErrCommentNotFound})
	}

	if err := h.service.Delete(c.Context(), commentIDUint64); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": domain.ErrCommentNotFound})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}
