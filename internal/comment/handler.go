package comment

import (
	"backend/internal/shared/ports/http"
	"context"
	"errors"
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

type LangMessage interface {
	GetResponseMessage(ctx context.Context, key string) string
}

type CommentHandler struct {
	validaotr http.Validator
	lang      LangMessage
	service   *CommentService
}

func NewCommentHandler(
	validator http.Validator,
	lang LangMessage,
	service *CommentService,
) *CommentHandler {
	return &CommentHandler{
		validaotr: validator,
		lang:      lang,
		service:   service,
	}
}

func (h *CommentHandler) Create(c *fiber.Ctx) error {
	comment, err := readBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	if validationErrors, statusCode, err := h.validaotr.ValidateStruct(c, comment); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.service.Create(c.Context(), comment); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), CreatedMessage),
		},
	)
}

func readBody(c *fiber.Ctx) (*CommentRequest, error) {
	comment := &CommentRequest{}
	if err := c.BodyParser(comment); err != nil {
		return nil, err
	}
	user_id, ok := c.Locals(UserIDKey).(uint64)
	if !ok {
		return nil, errors.New("user_id not found")
	}
	cardID, err := strconv.ParseUint(c.Params(CardIDKey), 10, 64)
	if err != nil {
		return nil, errors.New("card_id not found")
	}
	comment.UserID = user_id
	comment.CardID = cardID
	return comment, nil
}

func (h *CommentHandler) Update(c *fiber.Ctx) error {
	comment, err := readBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	commentID := c.Params(CommentIDKey)
	commentIDUint64, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "comment_id not found"})
	}
	comment.ID = commentIDUint64
	if validationErrors, statusCode, err := h.validaotr.ValidateStruct(c, comment); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.service.Update(c.Context(), comment); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}
	if err := h.service.Delete(c.Context(), commentIDUint64); err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}
