package eventbus

import (
	"fmt"
)

func NewEventBus[E any]() EventBus[E] {
	return EventBus[E]{
		handlers:     make(map[string]func(event E) error),
		errorHandler: func(event E, err error) {},
	}
}

type EventBus[E any] struct {
	handlers     map[string]func(event E) error
	errorHandler func(event E, err error)
}

func (e EventBus[E]) Emit(name string, event E) error {
	handler, exists := e.handlers[name]
	if !exists {
		return fmt.Errorf("There is no handler for event `%s`!", name)
	}

	err := handler(event)
	if err != nil && e.errorHandler != nil {
		e.errorHandler(event, err)
	}

	return err
}

func (e *EventBus[E]) On(name string, handler func(event E) error) {
	e.handlers[name] = handler
}

func (e *EventBus[E]) OnError(handler func(event E, err error)) {
	e.errorHandler = handler
}
