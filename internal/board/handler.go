package board

import "github.com/gofiber/fiber/v2"

type Validator interface {
	ValidateStruct(c *fiber.Ctx, structPtr interface{}) (map[string]string, int, error)
}

type BoardHandler struct {
	validator Validator
}

func NewBoardHandler(validator Validator) *BoardHandler {
	return &BoardHandler{
		validator: validator,
	}
}

func (h *BoardHandler) GetByUUID(c *fiber.Ctx) error   { return nil }
func (h *BoardHandler) GetList(c *fiber.Ctx) error     { return nil }
func (h *BoardHandler) GetProgress(c *fiber.Ctx) error { return nil }
func (h *BoardHandler) Store(c *fiber.Ctx) error       { return nil }
func (h *BoardHandler) Update(c *fiber.Ctx) error      { return nil }
func (h *BoardHandler) Delete(c *fiber.Ctx) error      { return nil }
