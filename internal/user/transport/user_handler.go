package transport

import (
	"backend/internal/shared/ports/http"
	"context"

	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey      = "userID"
	UpdatedMessage = "updated"
)

type UserService interface {
	Current(ctx context.Context, userID uint64) (*UserResponse, error)
	Update(ctx context.Context, userRequest *UserRequest, userID uint64) error
}

type UserHandler struct {
	validator   http.Validator
	lang        LangMessage
	userService UserService
}

func NewUserHandler(
	validator http.Validator,
	lang LangMessage,
	userService UserService,
) *UserHandler {
	return &UserHandler{
		validator:   validator,
		lang:        lang,
		userService: userService,
	}
}

func (h *UserHandler) Current(c *fiber.Ctx) error {
	userID, ok := c.Locals(UserIDKey).(uint64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userResponse, err := h.userService.Current(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.JSON(userResponse)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	body := &UserRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	userID, ok := c.Locals(UserIDKey).(uint64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	if err := h.userService.Update(c.Context(), body, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), UpdatedMessage),
		},
	)
}
