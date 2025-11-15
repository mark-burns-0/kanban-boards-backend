package transport

import (
	"backend/internal/board/domain"
	cardDomain "backend/internal/card/domain"
	boardError "backend/internal/shared/errors"
	"backend/internal/shared/ports/http"
	"backend/internal/shared/utils"
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	BoardIDKey  = "id"
	ColumnIDKey = "column_id"
	UserIDKey   = "userID"
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
	GetByUUID(ctx context.Context, board *domain.Board) (*domain.BoardWithDetails[cardDomain.CardWithComments], error)
	Create(ctx context.Context, board *domain.Board) error
	Update(ctx context.Context, board *domain.Board) error
	Delete(ctx context.Context, board *domain.Board) error
	CreateColumn(ctx context.Context, req *domain.BoardColumn) error
	UpdateColumn(ctx context.Context, req *domain.BoardColumn) error
	DeleteColumn(ctx context.Context, req *domain.BoardColumn) error
	MoveColumn(ctx context.Context, req *domain.BoardMoveCommand) error
}

type BoardHandler struct {
	validator    http.Validator
	lang         LangMessage
	boardService BoardService
	boardMapper  *BoardMapper
}

func NewBoardHandler(validator http.Validator, lang LangMessage, boardService BoardService) *BoardHandler {
	return &BoardHandler{
		validator:    validator,
		lang:         lang,
		boardService: boardService,
		boardMapper:  &BoardMapper{},
	}
}

func (h *BoardHandler) GetByUUID(c *fiber.Ctx) error {
	const op = "board.transport.handler.GetByUUID"
	body := &BoardRequest{}
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Missing board ID"})
	}
	body.ID = uuid
	userID, ok := c.Locals(UserIDKey).(uint64)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}
	body.UserID = userID

	response, err := h.boardService.GetByUUID(c.Context(), h.boardMapper.ToBoard(body))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBoardNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Board not found"})
		}
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Server error"})
	}

	return c.JSON(h.boardMapper.ToSingleBoardResponse(response))
}

func (h *BoardHandler) GetList(c *fiber.Ctx) error {
	const op = "board.transport.handler.GetList"
	body, err := utils.ParseBody[BoardGetFilter](c)
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

	response, err := h.boardService.GetList(c.Context(), h.boardMapper.ToBoardGetFilter(body))
	if err != nil {
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Server error"})
	}

	return c.JSON(h.boardMapper.ToBoardListResponse(response))
}

func (h *BoardHandler) Store(c *fiber.Ctx) error {
	const op = "board.transport.handler.Store"
	body, err := utils.ParseBody[BoardRequest](c)
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

	if err := h.boardService.Create(c.Context(), h.boardMapper.ToBoard(body)); err != nil {
		switch {
		case errors.Is(err, domain.ErrBoardAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"errors": "Board already exists"})
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

func (h *BoardHandler) Update(c *fiber.Ctx) error {
	const op = "board.transport.handler.Update"
	body, err := utils.ParseBody[BoardRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Missing board ID"})
	}
	body.ID = uuid

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

	if err := h.boardService.Update(c.Context(), h.boardMapper.ToBoard(body)); err != nil {
		switch {
		case errors.Is(err, domain.ErrBoardNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Board not found"})
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

func (h *BoardHandler) Delete(c *fiber.Ctx) error {
	const op = "board.transport.handler.Delete"
	body := &BoardRequest{}

	userID, ok := c.Locals(UserIDKey).(uint64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
	}
	body.UserID = userID

	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Missing board ID"})
	}
	body.ID = uuid

	if err := h.boardService.Delete(c.Context(), h.boardMapper.ToBoard(body)); err != nil {
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		switch {
		case errors.Is(err, boardError.ErrBoardHasCards):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"errors": boardError.ErrBoardHasCards.Error()})
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": domain.ErrBoardNotFound.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *BoardHandler) CreateColumn(c *fiber.Ctx) error {
	const op = "board.transport.handler.CreateColumn"
	body, err := utils.ParseBody[BoardColumnRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Missing board ID"})
	}
	body.BoardID = uuid

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

	if err := h.boardService.CreateColumn(c.Context(), h.boardMapper.ToBoardColumn(body)); err != nil {
		switch {
		case errors.Is(err, domain.ErrBoardNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Board not found"})
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

func (h *BoardHandler) UpdateColumn(c *fiber.Ctx) error {
	const op = "board.transport.handler.UpdateColumn"
	body, err := utils.ParseBody[BoardColumnRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Missing board ID"})
	}
	body.BoardID = uuid

	columnID := c.Params(ColumnIDKey)
	columnIDUint64, err := strconv.ParseUint(columnID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid column ID"})
	}
	body.ID = columnIDUint64

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

	if err := h.boardService.UpdateColumn(c.Context(), h.boardMapper.ToBoardColumn(body)); err != nil {
		switch {
		case errors.Is(err, domain.ErrColumnNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Column not found"})
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

func (h *BoardHandler) DeleteColumn(c *fiber.Ctx) error {
	const op = "board.transport.handler.DeleteColumn"
	body := &BoardColumnRequest{}
	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Missing board ID"})
	}
	columnID := c.Params(ColumnIDKey)

	columnIDUint64, err := strconv.ParseUint(columnID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid column ID"})
	}
	body.ID = columnIDUint64
	body.BoardID = uuid

	if err := h.boardService.DeleteColumn(c.Context(), h.boardMapper.ToBoardColumn(body)); err != nil {
		switch {
		case errors.Is(err, domain.ErrColumnNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Column not found"})
		case errors.Is(err, boardError.ErrColumnHasCards):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"errors": boardError.ErrColumnHasCards.Error()})
		}
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": domain.ErrColumnNotFound.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *BoardHandler) MoveColumn(c *fiber.Ctx) error {
	const op = "board.transport.handler.MoveColumn"
	body, err := utils.ParseBody[BoardColumnMoveRequest](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request body"})
	}

	uuid := c.Params(BoardIDKey)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Missing board ID"})
	}

	columnID := c.Params(ColumnIDKey)
	columnIDUint64, err := strconv.ParseUint(columnID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid column ID"})
	}

	body.BoardID = uuid
	body.ColumnID = columnIDUint64
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

	if err := h.boardService.MoveColumn(c.Context(), h.boardMapper.ToBoardMoveCommand(body)); err != nil {
		slog.Error(
			"service error",
			slog.String("operation", op),
			slog.Any("errors", err),
		)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": domain.ErrColumnNotFound.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": h.lang.GetResponseMessage(c.Context(), MovedMessage),
		},
	)
}
