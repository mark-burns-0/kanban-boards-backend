package events

type BoardDeletedEvent struct {
	BoardID string
	UserID  uint64
}

func (e BoardDeletedEvent) Name() string {
	return "BoardDeleted"
}

type ColumnDeletedEvent struct {
	ColumnID uint64
}

func (c ColumnDeletedEvent) Name() string {
	return "ColumnDeleted"
}
