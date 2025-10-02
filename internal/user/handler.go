package user

import (
	"backend/internal/platform/lang"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	validator    *validator.Validate
	langRegistry *lang.Registry
}

func NewAuthHandler(
	validator *validator.Validate,
	langRegistry *lang.Registry,
) *AuthHandler {
	return &AuthHandler{
		validator:    validator,
		langRegistry: langRegistry,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	body := UserRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	err := h.validator.Struct(body)
	if err != nil {
		validationMessages, err := h.langRegistry.Validate(c.Context(), err)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": validationMessages})
	}
	return c.JSON(fiber.Map{"data": body})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	body := UserRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	err := h.validator.Struct(body)
	if err != nil {
		validationMessages, err := h.langRegistry.Validate(c.Context(), err)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		if validationMessages != nil {
			return c.Status(fiber.StatusUnprocessableEntity).
				JSON(fiber.Map{"error": validationMessages})
		}
	}
	return c.JSON(fiber.Map{"data": body})
}

func (h *AuthHandler) Current(c *fiber.Ctx) error { return nil }

func (h *AuthHandler) Update(c *fiber.Ctx) error { return nil }

func (h *AuthHandler) Refresh(c *fiber.Ctx) error { return nil }
