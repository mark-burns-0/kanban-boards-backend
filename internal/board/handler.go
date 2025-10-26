package board

import (
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	BoardIDKey  = "id"
	ColumnIDKey = "column_id"
)

const (
	CreatedMessage = "created"
	UpdatedMessage = "updated"
	MovedMessage   = "moved"
)

type LangMessage interface {
	GetResponseMessage(ctx context.Context, key string) string
}

type BoardHandler struct {
	validator    http.Validator
	lang         LangMessage
	boardService *BoardService
}

func NewBoardHandler(validator http.Validator, lang LangMessage, boardService *BoardService) *BoardHandler {
	return &BoardHandler{
		validator:    validator,
		lang:         lang,
		boardService: boardService,
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

	response, err := h.boardService.GetList(c.Context(), body)
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

	if err := h.boardService.Create(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), CreatedMessage),
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

	if err := h.boardService.Update(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), UpdatedMessage),
		},
	)
}

func (h *BoardHandler) Delete(c *fiber.Ctx) error {
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	if err := h.boardService.Delete(c.Context(), uuid); err != nil {
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
	if err := h.boardService.CreateColumn(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), CreatedMessage),
		},
	)
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
	if err := h.boardService.UpdateColumn(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), UpdatedMessage),
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

	if err := h.boardService.DeleteColumn(c.Context(), body); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (h *BoardHandler) MoveToColumn(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardColumnMoveRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	columnID := c.Params(ColumnIDKey)
	columnIDUint64, err := strconv.ParseUint(columnID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid column ID"})
	}
	body.BoardID = uuid
	body.ColumnID = columnIDUint64

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.boardService.MoveToColumn(c.Context(), body); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), MovedMessage),
		},
	)
}
