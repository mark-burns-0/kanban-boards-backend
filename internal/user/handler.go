package user

import "github.com/gofiber/fiber/v2"

type Validator interface {
	ValidateStruct(c *fiber.Ctx, structPtr interface{}) error
}

type UserHandler struct {
	validator Validator
}

func NewUserHandler(validator Validator) *UserHandler {
	return &UserHandler{
		validator: validator,
	}
}

func (h *UserHandler) Current(c *fiber.Ctx) error { return nil }

func (h *UserHandler) Update(c *fiber.Ctx) error { return nil }

func (h *UserHandler) Refresh(c *fiber.Ctx) error { return nil }
