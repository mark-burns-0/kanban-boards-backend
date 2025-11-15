package transport

import (
	"backend/internal/card/domain"
	cardError "backend/internal/shared/errors"
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	CardIDKey  = "card_id"
	BoardIDKey = "id"
)

const (
	CreatedMessage = "created"
	UpdatedMessage = "updated"
	MovedMessage   = "moved"
)

type LangMessage interface {
	GetResponseMessage(ctx context.Context, key string) string
}

type CardService interface {
	GetListWithComments(ctx context.Context, boardID string) ([]*domain.CardWithComments, error)
	Create(ctx context.Context, req *domain.Card) error
	Update(ctx context.Context, req *domain.Card) error
	Delete(ctx context.Context, req *domain.Card) error
	MoveToNewPosition(ctx context.Context, req *domain.CardMoveCommand) error
}

type CardHandler struct {
	validator   http.Validator
	lang        LangMessage
	cardService CardService
	cardMapper  *CardMapper
}

func NewCardHandler(validator http.Validator, lang LangMessage, cardService CardService) *CardHandler {
	return &CardHandler{
		validator:   validator,
		lang:        lang,
		cardService: cardService,
		cardMapper:  &CardMapper{},
	}
}

func (h *CardHandler) Create(c *fiber.Ctx) error {
	const op = "card.transport.handler.Create"
	body, err := utils.ParseBody[CardRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	body.BoardID = c.Params(BoardIDKey)
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

	if err := h.cardService.Create(c.Context(), h.cardMapper.ToCard(body)); err != nil {
		switch {
		case errors.Is(err, domain.ErrCardAlreadyExists):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Card already exists"})
		}
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), CreatedMessage),
		},
	)
}

func (h *CardHandler) Update(c *fiber.Ctx) error {
	const op = "card.transport.handler.Update"
	body, err := utils.ParseBody[CardRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	cardID := c.Params(CardIDKey)
	cardIDUint64, err := strconv.ParseUint(cardID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	body.BoardID = c.Params(BoardIDKey)
	body.ID = cardIDUint64

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

	if err := h.cardService.Update(c.Context(), h.cardMapper.ToCard(body)); err != nil {
		switch {
		case errors.Is(err, domain.ErrCardNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Card not found"})
		}
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Server error"})
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), UpdatedMessage),
		},
	)
}

func (h *CardHandler) Delete(c *fiber.Ctx) error {
	const op = "card.transport.handler.Delete"
	body := &CardRequest{}

	cardID := c.Params(CardIDKey)
	cardIDUint64, err := strconv.ParseUint(cardID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	body.BoardID = c.Params(BoardIDKey)
	body.ID = cardIDUint64

	if err := h.cardService.Delete(c.Context(), h.cardMapper.ToCard(body)); err != nil {
		switch {
		case errors.Is(err, cardError.ErrCardHasComments):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"errors": cardError.ErrCardHasComments.Error(),
			})
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": domain.ErrCardNotFound.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CardHandler) MoveToNewPosition(c *fiber.Ctx) error {
	const op = "card.transport.handler.MoveToNewPosition"
	body, err := utils.ParseBody[CardMoveRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}
	cardID := c.Params(CardIDKey)
	cardIDUint64, err := strconv.ParseUint(cardID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	body.BoardID = c.Params(BoardIDKey)
	body.ID = cardIDUint64

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

	if err := h.cardService.MoveToNewPosition(c.Context(), h.cardMapper.ToCardMoveCommand(body)); err != nil {
		switch {
		case errors.Is(err, domain.ErrCardNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Card not found"})
		}
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Server error"})
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), MovedMessage),
		},
	)
}
