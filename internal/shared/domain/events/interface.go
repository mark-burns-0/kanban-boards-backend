package events

import "context"

type Event interface {
	Name() string
}

type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

type EventDispatcher interface {
	Dispatch(ctx context.Context, event Event) error
	Subscribe(eventName string, handler EventHandler) error
}
