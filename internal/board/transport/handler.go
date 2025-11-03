package transport

import (
	"backend/internal/board/domain"
	cardDomain "backend/internal/card/domain"
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

type BoardService interface {
	GetList(ctx context.Context, filter *domain.BoardGetFilter) (*domain.BoardListResult, error)
	GetByUUID(ctx context.Context, boardUUID string) (*domain.BoardWithDetails[cardDomain.CardWithComments], error)
	Create(ctx context.Context, board *domain.Board) error
	Update(ctx context.Context, board *domain.Board) error
	Delete(ctx context.Context, boardUUID string) error
	CreateColumn(ctx context.Context, req *domain.BoardColumn) error
	UpdateColumn(ctx context.Context, req *domain.BoardColumn) error
	DeleteColumn(ctx context.Context, req *domain.BoardColumn) error
	MoveColumn(ctx context.Context, req *domain.BoardMoveCommand) error
}

type BoardHandler struct {
	validator    http.Validator
	lang         LangMessage
	boardService BoardService
<<<<<<< Updated upstream
=======
	boardMapper  *BoardMapper
>>>>>>> Stashed changes
}

func NewBoardHandler(validator http.Validator, lang LangMessage, boardService BoardService) *BoardHandler {
	return &BoardHandler{
		validator:    validator,
		lang:         lang,
		boardService: boardService,
<<<<<<< Updated upstream
=======
		boardMapper:  &BoardMapper{},
>>>>>>> Stashed changes
	}
}

func (h *BoardHandler) GetByUUID(c *fiber.Ctx) error {
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
<<<<<<< Updated upstream

=======
>>>>>>> Stashed changes
	response, err := h.boardService.GetByUUID(c.Context(), uuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
<<<<<<< Updated upstream
	return c.JSON(response)
=======
	return c.JSON(h.boardMapper.ToSingleBoardResponse(response))
>>>>>>> Stashed changes
}

func (h *BoardHandler) GetList(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardGetFilter](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
<<<<<<< Updated upstream

=======
>>>>>>> Stashed changes
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
<<<<<<< Updated upstream

	response, err := h.boardService.GetList(c.Context(), body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return c.JSON(response)
=======
	response, err := h.boardService.GetList(c.Context(), h.boardMapper.ToBoardGetFilter(body))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return c.JSON(h.boardMapper.ToBoardListResponse(response))
>>>>>>> Stashed changes
}

func (h *BoardHandler) Store(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}

<<<<<<< Updated upstream
	if err := h.boardService.Create(c.Context(), body); err != nil {
=======
	if err := h.boardService.Create(c.Context(), h.boardMapper.ToBoard(body)); err != nil {
>>>>>>> Stashed changes
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	body.ID = uuid

	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}

<<<<<<< Updated upstream
	if err := h.boardService.Update(c.Context(), body); err != nil {
=======
	if err := h.boardService.Update(c.Context(), h.boardMapper.ToBoard(body)); err != nil {
>>>>>>> Stashed changes
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": domain.ErrBoardNotFound})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (h *BoardHandler) CreateColumn(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardColumnRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing board ID"})
	}
	body.BoardID = uuid
<<<<<<< Updated upstream

=======
>>>>>>> Stashed changes
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
<<<<<<< Updated upstream
	if err := h.boardService.CreateColumn(c.Context(), body); err != nil {
=======
	if err := h.boardService.CreateColumn(c.Context(), h.boardMapper.ToBoardColumn(body)); err != nil {
>>>>>>> Stashed changes
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
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
			return c.Status(statusCode).JSON(fiber.Map{"error": err})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
<<<<<<< Updated upstream
	if err := h.boardService.UpdateColumn(c.Context(), body); err != nil {
=======
	if err := h.boardService.UpdateColumn(c.Context(), h.boardMapper.ToBoardColumn(body)); err != nil {
>>>>>>> Stashed changes
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
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
<<<<<<< Updated upstream

	if err := h.boardService.DeleteColumn(c.Context(), body); err != nil {
=======
	if validationErrors, statusCode, err := h.validator.ValidateStruct(c, body); validationErrors != nil {
		if err != nil {
			return c.Status(statusCode).JSON(fiber.Map{"error": err})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
	if err := h.boardService.DeleteColumn(c.Context(), h.boardMapper.ToBoardColumn(body)); err != nil {
>>>>>>> Stashed changes
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (h *BoardHandler) MoveColumn(c *fiber.Ctx) error {
	body, err := utils.ParseBody[BoardColumnMoveRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
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
			return c.Status(statusCode).JSON(fiber.Map{"error": err})
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": validationErrors})
	}
<<<<<<< Updated upstream
	if err := h.boardService.MoveColumn(c.Context(), body); err != nil {
=======
	if err := h.boardService.MoveColumn(c.Context(), h.boardMapper.ToBoardMoveCommand(body)); err != nil {
>>>>>>> Stashed changes
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), MovedMessage),
		},
	)
}
