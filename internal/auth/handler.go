package auth

import (
	"github.com/gofiber/fiber/v2"
)

type Validator interface {
	ValidateStruct(c *fiber.Ctx, structPtr interface{}) error
}

type AuthHandler struct {
	validator   Validator
	authService *AuthService
}

func NewAuthHandler(
	validator Validator,
	authService *AuthService,
) *AuthHandler {
	return &AuthHandler{
		validator:   validator,
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	body := UserLoginRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.validator.ValidateStruct(c, body); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": body})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	body := UserCreateRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.validator.ValidateStruct(c, body); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": body})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	body := UserCreateRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.validator.ValidateStruct(c, body); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": body})
}
