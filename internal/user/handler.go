package user

import "github.com/gofiber/fiber/v2"

type Validator interface {
	ValidateStruct(c *fiber.Ctx, structPtr interface{}) error
}

type UserHandler struct {
	validator   Validator
	userService *UserService
}

func NewUserHandler(
	validator Validator,
	userService *UserService,
) *UserHandler {
	return &UserHandler{
		validator:   validator,
		userService: userService,
	}
}

func (h *UserHandler) Current(c *fiber.Ctx) error { return nil }

func (h *UserHandler) Update(c *fiber.Ctx) error { return nil }
