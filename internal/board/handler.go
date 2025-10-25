package board

import (
	"backend/internal/shared/ports/http"

	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey = "userID"
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
	body := &BoardGetFilter{}
	err := c.BodyParser(body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	userID, ok := c.Locals(UserIDKey).(uint64)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	body.UserID = userID

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
	body, err := bodyRead(c)
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
	uuid := c.Params("id")
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	body, err := bodyRead(c)
	body.ID = uuid
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
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
	uuid := c.Params("id")
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

func bodyRead(c *fiber.Ctx) (*BoardRequest, error) {
	body := &BoardRequest{}
	if err := c.BodyParser(body); err != nil {
		return nil, err
	}
	userID, ok := c.Locals(UserIDKey).(uint64)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}
	body.UserID = userID
	return body, nil
}
