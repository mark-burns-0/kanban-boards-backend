package comment

import (
	"backend/internal/shared/ports/http"

	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	validaotr http.Validator
}

func NewCommentHandler(validator http.Validator) *CommentHandler {
	return &CommentHandler{
		validaotr: validator,
	}
}

func (h *CommentHandler) Create(ctx *fiber.Ctx) error { return nil }
func (h *CommentHandler) Update(ctx *fiber.Ctx) error { return nil }
func (h *CommentHandler) Delete(ctx *fiber.Ctx) error { return nil }
