package eventbus

import (
	"context"
	"fmt"
)

func NewEventBus[E any]() EventBus[E] {
	return EventBus[E]{
		handlers:     make(map[string]Handler[E]),
		errorHandler: nil,
	}
}

type EventBus[E any] struct {
	handlers     map[string]Handler[E]
	errorHandler ErrorHandler[E]
}

type Handler[E any] func(ctx context.Context, event E) error

type ErrorHandler[E any] func(ctx context.Context, event E, err error)

func (e EventBus[E]) Emit(ctx context.Context, name string, event E) error {
	handler, exists := e.handlers[name]
	if !exists {
		return fmt.Errorf("There is no handler for event `%s`!", name)
	}

	err := handler(ctx, event)
	if err != nil && e.errorHandler != nil {
		e.errorHandler(ctx, event, err)
	}

	return err
}

func (e *EventBus[E]) On(name string, handler Handler[E]) {
	e.handlers[name] = handler
}

func (e *EventBus[E]) OnError(handler ErrorHandler[E]) {
	e.errorHandler = handler
}
