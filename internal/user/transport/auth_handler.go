package transport

import (
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"
	"context"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	RefreshTokenHeader = "Refresh-Token"
	BearerPrefix       = "Bearer "
)

type AuthService interface {
	Register(ctx context.Context, userRequest *UserRegisterRequest) error
	Login(ctx context.Context, userRequest *UserLoginRequest) (*TokensResponse, error)
	RefreshToken(ctx context.Context, token string) (*TokensResponse, error)
}

type AuthHandler struct {
	validator   http.Validator
	lang        http.LangMessage
	authService AuthService
}

func NewAuthHandler(
	validator http.Validator,
	lang http.LangMessage,
	authService AuthService,
) *AuthHandler {
	return &AuthHandler{
		validator:   validator,
		lang:        lang,
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	body, err := utils.ParseBody[UserLoginRequest](c)
	if err != nil {
		slog.Error("Invalid request body", slog.Any("err", err))
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid request body"})
	}
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	tokenResponse, err := h.authService.Login(c.Context(), body)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Server error"})
	}
	return c.JSON(tokenResponse)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	const op = "user.transport.auth_handler.Register"
	body, err := utils.ParseBody[UserRegisterRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid request body"})
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
	if err := h.authService.Register(c.Context(), body); err != nil {
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("error", err),
		)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Server error"})
	}
	return c.JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	token := c.Get(RefreshTokenHeader)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing refresh token"})
	}
	token = strings.TrimPrefix(token, BearerPrefix)
	tokenResponse, err := h.authService.RefreshToken(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.JSON(tokenResponse)
}
