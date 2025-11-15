package events

import (
	"backend/internal/shared/domain/events"
	"context"
	"sync"
)

type inMemoryBus struct {
	mu            sync.RWMutex
	subscriptions map[string][]events.EventHandler
}

func NewInMemoryBus() *inMemoryBus {
	return &inMemoryBus{
		subscriptions: make(map[string][]events.EventHandler, 10),
	}
}

func (ib *inMemoryBus) Dispatch(ctx context.Context, event events.Event) error {
	ib.mu.RLock()
	handlers, ok := ib.subscriptions[event.Name()]
	ib.mu.RUnlock()
	if !ok {
		return nil
	}
	for _, h := range handlers {
		if err := h.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (ib *inMemoryBus) Subscribe(eventName string, handler events.EventHandler) error {
	ib.mu.Lock()
	defer ib.mu.Unlock()
	ib.subscriptions[eventName] = append(ib.subscriptions[eventName], handler)
	return nil
}
