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

func (r *CardRepository) Create(ctx context.Context, card *Card) error {
	const op = "card.repository.Create"
	query := `
		INSERT INTO cards (board_id, column_id, text, description, position, properties)
		VALUES($1, $2, $3, $4, $5, $6)
	`
	return utils.OpExec(
		ctx,
		r.storage,
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
		SET column_id = $1, position = $2, text = $3, description = $4, properties = $5, updated_at = NOW()
		WHERE id = $6 AND board_id = $7 AND deleted_at IS NULL
	`
	return utils.OpExec(
		ctx,
		r.storage,
		op,
		query,
		ErrCardNotFound,
		card.ColumnID,
		card.Position,
		card.Text,
		card.Description,
		card.cardProperties,
		card.ID,
		card.BoardID,
	)
}

func (r *CardRepository) Delete(ctx context.Context, card *Card) error {
	const op = "card.repository.Delete"
	var exists bool

	row := r.storage.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT 1 FROM cards WHERE id = $1 AND board_id = $2)",
		card.ID,
		card.BoardID,
	)
	err := row.Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrCardNotFound)
	}
	query := `UPDATE cards SET deleted_at = NOW() WHERE id = $1 AND board_id = $2`
	result, err := r.storage.ExecContext(ctx, query, card.ID, card.BoardID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("%s: %w", op, ErrCardNotFound)
	}
	return nil
}

func (r *CardRepository) MoveToNewPosition(
	ctx context.Context, boardID string, cardID, fromColumnID, toColumnID, cardPosition uint64,
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
		r.storage,
		op,
		query,
		ErrCardNotFound,
		toColumnID,
		cardPosition,
		cardID,
		boardID,
	)
	if err != nil {
		return err
	}
	return nil
}
