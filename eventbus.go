package eventbus

import (
	"context"
	"fmt"
	"sync"
)

func NewEventBus[E any]() EventBus[E] {
	return EventBus[E]{
		handlers:     make(map[string][]Handler[E]),
		errorHandler: nil,
	}
}

type EventBus[E any] struct {
	handlers     map[string][]Handler[E]
	errorHandler ErrorHandler[E]
	mux          sync.Mutex
}

type Handler[E any] func(ctx context.Context, event E) error

type ErrorHandler[E any] func(ctx context.Context, event E, err error)

func (e *EventBus[E]) Emit(ctx context.Context, name string, event E) (err error) {
	e.mux.Lock()
	handlers, exists := e.handlers[name]
	e.mux.Unlock()

	if !exists {
		return fmt.Errorf("there is no handler for event `%s`", name)
	}

	for _, handler := range handlers {
		if err = handler(ctx, event); err != nil {
			if e.errorHandler != nil {
				e.errorHandler(ctx, event, err)
			}

			break
		}
	}

	return err
}

func (e *EventBus[E]) Subscribe(name string, handler Handler[E]) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.handlers[name] = append(e.handlers[name], handler)
}

func (e *EventBus[E]) OnError(handler ErrorHandler[E]) {
	e.errorHandler = handler
}
