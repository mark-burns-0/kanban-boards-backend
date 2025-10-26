package card

import (
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"
	"context"
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

type CardHandler struct {
	validator   http.Validator
	lang        LangMessage
	cardService *CardService
}

func NewCardHandler(validator http.Validator, lang LangMessage, cardService *CardService) *CardHandler {
	return &CardHandler{
		validator:   validator,
		lang:        lang,
		cardService: cardService,
	}
}

func (h *CardHandler) Create(c *fiber.Ctx) error {
	body, err := utils.ParseBody[CardRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	body.BoardID = c.Params(BoardIDKey)
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.cardService.Create(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), CreatedMessage),
		},
	)
}

func (h *CardHandler) Update(c *fiber.Ctx) error {
	body, err := utils.ParseBody[CardRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	cardID := c.Params(CardIDKey)
	cardIDUint64, err := strconv.ParseUint(cardID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	body.BoardID = c.Params(BoardIDKey)
	body.ID = cardIDUint64
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.cardService.Update(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), UpdatedMessage),
		},
	)
}

func (h *CardHandler) Delete(c *fiber.Ctx) error {
	body := &CardRequest{}
	cardID := c.Params(CardIDKey)
	cardIDUint64, err := strconv.ParseUint(cardID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	body.BoardID = c.Params(BoardIDKey)
	body.ID = cardIDUint64
	if err := h.cardService.Delete(c.Context(), body); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": ErrCardNotFound.Error(),
		})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (h *CardHandler) MoveToNewPosition(c *fiber.Ctx) error {
	body, err := utils.ParseBody[CardMoveRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	body.BoardID = c.Params(BoardIDKey)
	cardID := c.Params(CardIDKey)
	cardIDUint64, err := strconv.ParseUint(cardID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	body.ID = cardIDUint64
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.cardService.MoveToNewPosition(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), MovedMessage),
		},
	)
}
