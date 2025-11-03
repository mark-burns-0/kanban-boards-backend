package transport

import (
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"
	"backend/internal/user/domain"
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey      = "userID"
	UpdatedMessage = "updated"
)

type UserService interface {
	Current(ctx context.Context, userID uint64) (*domain.User, error)
	Update(ctx context.Context, req *domain.User, userID uint64) error
}

type UserHandler struct {
	validator   http.Validator
	lang        http.LangMessage
	userService UserService
	mapperUser  *UserMapper
}

func NewUserHandler(
	validator http.Validator,
	lang http.LangMessage,
	userService UserService,
) *UserHandler {
	return &UserHandler{
		validator:   validator,
		lang:        lang,
		userService: userService,
		mapperUser:  &UserMapper{},
	}
}

func (h *UserHandler) Current(c *fiber.Ctx) error {
	const op = "user.transport.user_handler.Current"
	userID, ok := c.Locals(UserIDKey).(uint64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userResponse, err := h.userService.Current(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.JSON(h.mapperUser.ToUserResponse(userResponse))
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	const op = "user.transport.user_handler.Update"
	body, err := utils.ParseBody[UserRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			slog.Error("validator error",
				slog.String("op", op),
				slog.Any("err", err),
			)
			return c.Status(statusCode).JSON(fiber.Map{"error": "Validation error"})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}

	userID, ok := c.Locals(UserIDKey).(uint64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.userService.Update(c.Context(), h.mapperUser.ToUserDomain(body), userID); err != nil {
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("error", err),
		)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), UpdatedMessage),
		},
	)
}
