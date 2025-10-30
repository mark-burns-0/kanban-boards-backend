package card

import (
	"backend/internal/shared/dto"
	"cmp"
	"context"
	"fmt"
	"slices"
)

type CardGetter interface {
	GetListWithComments(ctx context.Context, boardID string) ([]*CardWithComments, error)
	GetMaxColumnPosition(ctx context.Context, boardUUID string, columnID uint64) (uint64, error)
	GetById(ctx context.Context, card *Card) (*Card, error)
}

type CardCreator interface {
	Create(context.Context, *Card) error
	Exists(ctx context.Context, card *Card) (bool, error)
}

type CardUpdater interface {
	Update(context.Context, *Card) error
	MoveToNewPosition(ctx context.Context, boardID string, cardID, fromColumnID, toColumnID, cardFromPosition, cardToPosition uint64) error
}

type CardDeleter interface {
	Delete(context.Context, *Card) error
}

type CardRepo interface {
	CardGetter
	CardCreator
	CardUpdater
	CardDeleter
}

type CardService struct {
	repo CardRepo
}

func NewCardService(repo CardRepo) *CardService {
	return &CardService{
		repo: repo,
	}
}

func (s *CardService) GetListWithComments(ctx context.Context, boardID string) ([]*dto.CardWithComments, error) {
	const op = "card.service.GetListWithComments"
	raws, err := s.repo.GetListWithComments(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	response := make([]*dto.CardWithComments, 0, len(raws))

	for _, raw := range raws {
		card := &dto.CardWithComments{
			ID:          raw.ID,
			ColumnID:    raw.ColumnID,
			Position:    raw.Position,
			BoardID:     raw.BoardID,
			Text:        raw.Text,
			Description: raw.Description,
			CreatedAt:   raw.CreatedAt,
			Properties: &dto.CardProperties{
				Color: raw.cardProperties.Color,
				Tag:   raw.cardProperties.Tag,
			},
			Comments: make([]*dto.CardComment, 0, len(raw.Comments)),
		}
		var comment *dto.CardComment
		for _, rawComment := range raw.Comments {
			comment = &dto.CardComment{
				ID:        rawComment.ID,
				CardID:    rawComment.CardID,
				Text:      rawComment.Text,
				CreatedAt: rawComment.CreatedAt,
			}
			card.Comments = append(card.Comments, comment)
		}

		response = append(response, card)
	}
	slices.SortFunc(response, func(a, b *dto.CardWithComments) int {
		return cmp.Compare(*a.ID, *b.ID)
	})
	return response, nil
}

func (s *CardService) Create(ctx context.Context, req *CardRequest) error {
	const op = "card.service.Create"
	maxPosition, err := s.repo.GetMaxColumnPosition(ctx, req.BoardID, req.ColumnID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	card := &Card{
		BoardID:     req.BoardID,
		ColumnID:    req.ColumnID,
		Text:        req.Text,
		Position:    maxPosition + 1,
		Description: req.Description,
		cardProperties: cardProperties{
			Color: req.Color,
			Tag:   req.Tag,
		},
	}
	if err := s.repo.Create(ctx, card); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *CardService) Update(ctx context.Context, req *CardRequest) error {
	const op = "card.service.Update"
	card := &Card{
		ID:          req.ID,
		ColumnID:    req.ColumnID,
		BoardID:     req.BoardID,
		Text:        req.Text,
		Description: req.Description,
		cardProperties: cardProperties{
			Color: req.Color,
			Tag:   req.Tag,
		},
	}
	exists, err := s.repo.Exists(ctx, card)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrCardNotFound)
	}
	if err := s.repo.Update(ctx, card); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *CardService) Delete(ctx context.Context, req *CardRequest) error {
	const op = "card.service.Delete"
	card := &Card{
		ID:      req.ID,
		BoardID: req.BoardID,
	}
	exists, err := s.repo.Exists(ctx, card)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrCardNotFound)
	}
	card, err = s.repo.GetById(ctx, card)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := s.repo.Delete(ctx, card); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Описание алгоритма MoveCard
// Общая структура
// Метод MoveCard в репозитории принимает параметры перемещения карточки: идентификаторы доски, карточки, исходной и целевой колонок, исходной и целевой позиций. Весь алгоритм выполняется в транзакции с уровнем изоляции Repeatable Read.
// Валидация и проверки
// Первым делом проверяется случай, когда исходная и целевая колонки совпадают и исходная позиция равна целевой. В этом случае метод завершается без выполнения операций.
// Если колонки разные, выполняется перемещение между колонками. Если колонки одинаковые, но позиции разные, выполняется перемещение внутри колонки.
// Перемещение внутри колонки
// Для перемещения внутри одной колонки используется стратегия временных позиций. Сначала перемещаемая карточка устанавливается на временную позицию, вычисляемую как отрицательное значение исходной позиции. Это безопасно, так как в системе позиции всегда положительные.
// Затем, в зависимости от направления перемещения, выполняется сдвиг других карточек. Если карточка перемещается влево (на меньшую позицию), то все карточки между целевой и исходной позицией сдвигаются вправо. Если карточка перемещается вправо (на большую позицию), то карточки между исходной и целевой позицией сдвигаются влево.
// После сдвига карточка устанавливается на целевую позицию с проверкой по временной позиции.
// Перемещение между колонками
// При перемещении между разными колонками также используется стратегия временных позиций. Перемещаемая карточка сначала устанавливается на временную позицию в исходной колонке.
// Затем в исходной колонке все карточки, следующие после исходной позиции, сдвигаются на одну позицию влево для заполнения освободившегося места.
// В целевой колонке все карточки, находящиеся на целевую позицию и далее, сдвигаются на одну позицию вправо для освобождения места под новую карточку.
// Наконец, карточка устанавливается на целевую позицию в целевой колонке с проверкой по временной позиции.
// Обработка особых случаев
// При перемещении в пустую колонку карточка устанавливается на первую позицию без сдвига в целевой колонке.
// Если целевая позиция превышает максимально допустимую в колонке, возвращается ошибка.
// Все операции выполняются атомарно в рамках одной транзакции, что гарантирует целостность данных при параллельных операциях.
func (s *CardService) MoveToNewPosition(ctx context.Context, req *CardMoveRequest) error {
	const op = "card.service.MoveToNewPosition"
	card := &Card{
		ID:      req.ID,
		BoardID: req.BoardID,
	}
	exists, err := s.repo.Exists(ctx, card)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, ErrCardNotFound)
	}
	if err := s.repo.MoveToNewPosition(
		ctx,
		req.BoardID,
		req.ID,
		req.FromColumnID,
		req.ToColumnID,
		req.FromPosition,
		req.ToPosition,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
