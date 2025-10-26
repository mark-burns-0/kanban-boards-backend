package board

import (
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey  = "userID"
	BoardIDKey = "id"
)

type BoardHandler struct {
	validator http.Validator
	service   *BoardService
}

func NewBoardHandler(validator http.Validator, service *BoardService) *BoardHandler {
	return &BoardHandler{
		validator: validator,
		service:   service,
	}
}

func (h *BoardHandler) GetByUUID(c *fiber.Ctx) error { return nil }

func (h *BoardHandler) GetList(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardGetFilter](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}

	response, err := h.service.GetList(c.Context(), body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(response)
}

func (h *BoardHandler) Store(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}

	if err := h.service.Create(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": "Board created successfully",
		},
	)
}

func (h *BoardHandler) Update(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	body.ID = uuid

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}

	if err := h.service.Update(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Board updated successfully"})
}

func (h *BoardHandler) Delete(c *fiber.Ctx) error {
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	if err := h.service.Delete(c.Context(), uuid); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": ErrBoardNotFound.Error()})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (h *BoardHandler) CreateColumn(*fiber.Ctx) error {
	return nil
}

func (h *BoardHandler) UpdateColumn(*fiber.Ctx) error {
	return nil
}

func (h *BoardHandler) DeleteColumn(*fiber.Ctx) error {
	return nil
}

func (h *BoardHandler) MoveToColumn(c *fiber.Ctx) error {
	return nil
}
