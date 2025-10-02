package eventbus

import (
	"context"
	"sync"
)

func newEvent[Payload any]() *Event[Payload] {
	return &Event[Payload]{
		handlers: make([]Handler[Payload], 0),
	}
}

type Event[Payload any] struct {
	handlers []Handler[Payload]
	mux      sync.Mutex
}

type Handler[Payload any] func(ctx context.Context, p Payload) error

func (e *Event[Payload]) Emit(ctx context.Context, p Payload) (err error) {
	e.mux.Lock()
	handlers := e.handlers
	e.mux.Unlock()

	for _, handler := range handlers {
		if err = handler(ctx, p); err != nil {
			break
		}
	}

	return err
}

func (e *Event[Payload]) Subscribe(handler Handler[Payload]) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.handlers = append(e.handlers, handler)
}
