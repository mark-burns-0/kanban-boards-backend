package card

import (
	"backend/internal/shared/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCardAlreadyExists = errors.New("card already exists")
	ErrCardNotFound      = errors.New("card not found")
)

type Storage interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)

	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type CardRepository struct {
	storage Storage
}

func NewCardRepository(
	storage Storage,
) *CardRepository {
	return &CardRepository{
		storage: storage,
	}
}

func (r *CardRepository) GetListWithComments(ctx context.Context, boardID string) ([]*CardWithComments, error) {
	const op = "card.repository.GetListWithComment"
	query := `
		SELECT
				cards.id,
				cards.board_id,
				cards.column_id,
				cards.text,
				cards.description,
				cards.position,
				cards.properties,
				cards.created_at,
				comments.id,
				comments.card_id,
				comments.text,
				comments.created_at
		FROM cards
		LEFT JOIN comments ON comments.card_id = cards.id
		WHERE cards.deleted_at is null and comments.deleted_at is NULL
			and cards.board_id = $1
	`
	rows, err := r.storage.QueryContext(ctx, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	cardComments := make(map[uint64]*CardWithComments)
	for rows.Next() {
		var cardID, cardPosition, columnID, commentID, commentCardID *uint64
		var boardID, cardText, cardDescription, commentText *string
		var properties *cardProperties
		var cardCreatedAt, commentCreatedAt *time.Time

		err := rows.Scan(
			&cardID, &boardID, &columnID, &cardText, &cardDescription, &cardPosition, &properties, &cardCreatedAt,
			&commentID, &commentCardID, &commentText, &commentCreatedAt,
		)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if cardID == nil {
			continue
		}
		_, ok := cardComments[*cardID]
		if !ok {
			cardComments[*cardID] = &CardWithComments{
				ID:             cardID,
				BoardID:        boardID,
				ColumnID:       columnID,
				Text:           cardText,
				Description:    cardDescription,
				Position:       cardPosition,
				cardProperties: properties,
				CreatedAt:      cardCreatedAt,
				Comments:       []*CardComment{},
			}
		}
		if commentID != nil {
			cardComments[*cardID].Comments = append(cardComments[*cardID].Comments, &CardComment{
				ID:        commentID,
				CardID:    commentCardID,
				Text:      commentText,
				CreatedAt: commentCreatedAt,
			})
		}
	}
	cardWithComments := make([]*CardWithComments, 0, len(cardComments))
	for _, card := range cardComments {
		cardWithComments = append(cardWithComments, card)
	}
	return cardWithComments, nil
}

func (r *CardRepository) GetById(ctx context.Context, card *Card) (*Card, error) {
	const op = "card.repository.GetById"
	data := &Card{}
	query := "SELECT id, column_id, board_id, text, description, position FROM cards WHERE id = $1 AND deleted_at IS NULL"
	row := r.storage.QueryRowContext(ctx, query, card.ID)
	err := row.Scan(
		&data.ID,
		&data.ColumnID,
		&data.BoardID,
		&data.Text,
		&data.Description,
		&data.Position,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrCardNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return data, nil
}

func (r *CardRepository) Create(ctx context.Context, card *Card) error {
	const op = "card.repository.Create"
	query := `
		INSERT INTO cards (board_id, column_id, text, description, position, properties)
		VALUES($1, $2, $3, $4, $5, $6)
	`
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		ErrCardAlreadyExists,
		card.BoardID,
		card.ColumnID,
		card.Text,
		card.Description,
		card.Position,
		card.cardProperties,
	)
}

func (r *CardRepository) Update(ctx context.Context, card *Card) error {
	const op = "card.repository.Update"
	query := `
		UPDATE cards
		SET text = $1, description = $2, properties = $3, updated_at = NOW()
		WHERE id = $4 AND board_id = $5 AND deleted_at IS NULL
	`
	return utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		ErrCardNotFound,
		card.Text,
		card.Description,
		card.cardProperties,
		card.ID,
		card.BoardID,
	)
}

func (r *CardRepository) Delete(ctx context.Context, card *Card) error {
	const op = "card.repository.Delete"
	tx, err := r.storage.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	query := `UPDATE cards SET position = position - 1 WHERE column_id = $1 AND  board_id = $2 AND position > $3 AND deleted_at IS NULL`
	_, err = tx.ExecContext(ctx, query, card.ColumnID, card.BoardID, card.Position)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	query = `UPDATE cards SET deleted_at = NOW(), position = NULL WHERE id = $1 AND board_id = $2 AND deleted_at IS NULL`
	err = utils.OpExec(ctx, tx.ExecContext, op, query, ErrCardNotFound, card.ID, card.BoardID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return tx.Commit()
}

func (r *CardRepository) Exists(ctx context.Context, card *Card) (bool, error) {
	const op = "card.repository.Exists"
	row := r.storage.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT 1 FROM cards WHERE id = $1 AND board_id = $2)",
		card.ID,
		card.BoardID,
	)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *CardRepository) GetMaxColumnPosition(ctx context.Context, boardUUID string, columnID uint64) (uint64, error) {
	const op = "card.repository.GetMaxColumnPosition"
	var maxValue sql.NullInt64
	query := "SELECT MAX(position) FROM cards WHERE board_id = $1 AND column_id = $2 AND deleted_at IS NULL"
	row := r.storage.QueryRowContext(ctx, query, boardUUID, columnID)
	err := row.Scan(&maxValue)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if !maxValue.Valid {
		return 0, nil
	}
	return uint64(maxValue.Int64), nil
}

func (r *CardRepository) MoveToNewPosition(
	ctx context.Context, boardID string, cardID, fromColumnID, toColumnID, cardFromPosition, cardToPosition uint64,
) error {
	const op = "card.repository.MoveToNewPosition"
	if fromColumnID == toColumnID {
		return nil
	}
	query := `
		UPDATE cards SET column_id = $1, position = $2, updated_at = NOW()
		WHERE id = $3 AND board_id = $4 AND deleted_at IS NULL
	`
	err := utils.OpExec(
		ctx,
		r.storage.ExecContext,
		op,
		query,
		ErrCardNotFound,
		toColumnID,
		cardFromPosition,
		cardID,
		boardID,
	)
	if err != nil {
		return err
	}
	return nil
}
