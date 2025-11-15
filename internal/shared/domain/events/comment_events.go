package events

type CardDeletedEvent struct {
	CardID uint64
}

func (e CardDeletedEvent) Name() string {
	return "CardDeleted"
}
