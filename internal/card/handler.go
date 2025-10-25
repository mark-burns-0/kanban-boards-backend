package card

import (
	"backend/internal/shared/ports/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	CardIDKey  = "card_id"
	BoardIDKey = "id"
)

type CardHandler struct {
	validator   http.Validator
	cardService *CardService
}

func NewCardHandler(validator http.Validator, cardService *CardService) *CardHandler {
	return &CardHandler{
		validator:   validator,
		cardService: cardService,
	}
}

func (h *CardHandler) Create(c *fiber.Ctx) error {
	body, err := bodyRead(c)
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
			"message": "Card created successfully",
		},
	)
}

func (h *CardHandler) Update(c *fiber.Ctx) error {
	body, err := bodyRead(c)
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Card updated successfully",
	})
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

func bodyRead(ctx *fiber.Ctx) (*CardRequest, error) {
	board := &CardRequest{}
	if err := ctx.BodyParser(board); err != nil {
		return nil, err
	}
	return board, nil
}

func (h *CardHandler) MoveToNewPosition(c *fiber.Ctx) error {
	body := &CardMoveRequest{}
	if err := c.BodyParser(body); err != nil {
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Card moved successfully",
	})
}
