package board

import (
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/k0kubun/pp"
)

const (
	BoardIDKey  = "id"
	ColumnIDKey = "column_id"
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

func (h *BoardHandler) CreateColumn(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardColumnRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	body.BoardID = uuid

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.service.CreateColumn(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (h *BoardHandler) UpdateColumn(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardColumnRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	body.BoardID = uuid
	columnID := c.Params(ColumnIDKey)
	columnIDUint64, err := strconv.ParseUint(columnID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid column ID"})
	}
	body.ID = columnIDUint64
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.service.UpdateColumn(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "Column updated successfully",
		},
	)
}

func (h *BoardHandler) DeleteColumn(c *fiber.Ctx) error {
	body := &BoardColumnRequest{}
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	columnID := c.Params(ColumnIDKey)
	columnIDUint64, err := strconv.ParseUint(columnID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid column ID"})
	}
	body.ID = columnIDUint64
	body.BoardID = uuid

	if err := h.service.DeleteColumn(c.Context(), body); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (h *BoardHandler) MoveToColumn(c *fiber.Ctx) error {
	pp.Println(c)
	return nil
}
