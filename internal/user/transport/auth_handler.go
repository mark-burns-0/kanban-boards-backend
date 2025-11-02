package transport

import (
	"backend/internal/shared/ports/http"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	RefreshTokenHeader = "Refresh-Token"
	BearerPrefix       = "Bearer "
)

type LangMessage interface {
	GetResponseMessage(ctx context.Context, key string) string
}

type AuthService interface {
	Register(ctx context.Context, userRequest *UserCreateRequest) error
	Login(ctx context.Context, userRequest *UserLoginRequest) (*TokensResponse, error)
	RefreshToken(ctx context.Context, token string) (*TokensResponse, error)
}

type AuthHandler struct {
	validator   http.Validator
	lang        LangMessage
	authService AuthService
}

func NewAuthHandler(
	validator http.Validator,
	lang LangMessage,
	authService AuthService,
) *AuthHandler {
	return &AuthHandler{
		validator:   validator,
		lang:        lang,
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	body := &UserLoginRequest{}
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
	tokenResponse, err := h.authService.Login(c.Context(), body)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tokenResponse)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	body := &UserCreateRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.authService.Register(c.Context(), body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
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
