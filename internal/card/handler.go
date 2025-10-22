package card

import (
	"backend/internal/shared/ports/http"

	"github.com/gofiber/fiber/v2"
)

type CardHandler struct {
	validator http.Validator
}

func NewCardHandler(validator http.Validator) *CardHandler {
	return &CardHandler{
		validator: validator,
	}
}

func (h *CardHandler) Create(ctx *fiber.Ctx) error {
	return nil
}

func (h *CardHandler) Delete(ctx *fiber.Ctx) error {
	return nil
}

func (h *CardHandler) MoveToNewPosition(ctx *fiber.Ctx) error {
	return nil
}
