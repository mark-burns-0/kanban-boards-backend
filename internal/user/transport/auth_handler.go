package transport

import (
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"
	"backend/internal/user/domain"
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	BearerPrefix    = "Bearer "
	RefreshTokenKey = "refreshToken"
	CreatedMessage  = "created"
)

type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterCommand) error
	Login(ctx context.Context, req *domain.LoginCommand) (*domain.Tokens, error)
	RefreshToken(ctx context.Context, token string) (*domain.Tokens, error)
}

type AuthHandler struct {
	validator   http.Validator
	lang        http.LangMessage
	authService AuthService
	authMapper  AuthMapper
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
		authMapper:  AuthMapper{},
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	const op = "user.transport.auth_handler.Login"
	body, err := utils.ParseBody[UserLoginRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			slog.Error("validator error",
				slog.String("op", op),
				slog.Any("err", err),
			)
			return c.Status(statusCode).JSON(fiber.Map{"errors": "Validation error"})
		}
		return c.Status(statusCode).JSON(fiber.Map{"errors": validationErrors})
	}

	tokenResponse, err := h.authService.Login(c.Context(), h.authMapper.ToLoginCommand(body))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Email not found"})
		case errors.Is(err, domain.ErrInvalidPassword):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid password"})
		}
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": "Server error"})
	}

	return c.JSON(h.authMapper.ToResponseTokens(tokenResponse))
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	const op = "user.transport.auth_handler.Register"
	body, err := utils.ParseBody[UserRegisterRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"errors": "Invalid request body"})
	}

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			slog.Error("validator error",
				slog.String("op", op),
				slog.Any("err", err),
			)
			return c.Status(statusCode).JSON(fiber.Map{"errors": "Validation error"})
		}
		return c.Status(statusCode).JSON(fiber.Map{"errors": validationErrors})
	}

	if err := h.authService.Register(c.Context(), h.authMapper.ToRegisterCommand(body)); err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"errors": "User already exists"})
		case errors.Is(err, domain.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Email not found"})
		case errors.Is(err, domain.ErrInvalidPassword):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid password"})
		}
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": "Server error"})
	}

	return c.JSON(fiber.Map{
		"message": h.lang.GetResponseMessage(c.Context(), CreatedMessage),
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	const op = "user.transport.auth_handler.RefreshToken"
	token, ok := c.Locals(RefreshTokenKey).(string)
	if token == "" || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Unauthorized"})
	}

	token = strings.TrimPrefix(token, BearerPrefix)
	tokenResponse, err := h.authService.RefreshToken(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Unauthorized"})
	}

	return c.JSON(h.authMapper.ToResponseTokens(tokenResponse))
}
