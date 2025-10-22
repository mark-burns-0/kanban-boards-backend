package board

import (
	"backend/internal/shared/ports/http"

	"github.com/gofiber/fiber/v2"
)

type BoardHandler struct {
	validator http.Validator
}

func NewBoardHandler(validator http.Validator) *BoardHandler {
	return &BoardHandler{
		validator: validator,
	}
}

func (h *BoardHandler) GetByUUID(c *fiber.Ctx) error    { return nil }
func (h *BoardHandler) GetList(c *fiber.Ctx) error      { return nil }
func (h *BoardHandler) Store(c *fiber.Ctx) error        { return nil }
func (h *BoardHandler) Update(c *fiber.Ctx) error       { return nil }
func (h *BoardHandler) Delete(c *fiber.Ctx) error       { return nil }
func (h *BoardHandler) MoveToColumn(c *fiber.Ctx) error { return nil }
