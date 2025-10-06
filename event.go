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
	handlers   []Handler[Payload]
	ErrorEvent *ErrorEvent[Payload]
	mux        sync.Mutex
}

// Handler defines the function signature for event subscribers.
type Handler[Payload any] func(ctx context.Context, p Payload) error

// Trigger all subscribed handlers with the given payload.
//
// If any handler returns an error, execution stops and the error is returned.
// If an ErrorEvent is set, it will be emitted with the error and payload.
func (e *Event[Payload]) Emit(ctx context.Context, p Payload) (err error) {
	for _, handler := range e.cloneHandlers() {
		if err = handler(ctx, p); err != nil {
			if e.ErrorEvent != nil {
				e.ErrorEvent.Emit(ctx, err, p)
			}

			break
		}
	}

	return err
}

// Triggers all subscribed handlers concurrently with the given payload.
//
// Handlers run in separate goroutines, and any errors are sent to the ErrorEvent if set.
// Unlike Emit, this does not stop on the first error and waits for all handlers to complete.
func (e *Event[Payload]) EmitParallel(ctx context.Context, p Payload) {
	handlers := e.cloneHandlers()

	var wg sync.WaitGroup
	wg.Add(len(handlers))

	for _, handler := range handlers {
		go func(h Handler[Payload]) {
			defer wg.Done()

			if err := h(ctx, p); err != nil && e.ErrorEvent != nil {
				e.ErrorEvent.Emit(ctx, err, p)
			}
		}(handler)
	}

	wg.Wait()
}

func (e *Event[Payload]) cloneHandlers() []Handler[Payload] {
	e.mux.Lock()
	defer e.mux.Unlock()

	handlers := make([]Handler[Payload], len(e.handlers))
	copy(handlers, e.handlers)

	return handlers
}

// Register a new handler for the event.
func (e *Event[Payload]) Subscribe(handler Handler[Payload]) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.handlers = append(e.handlers, handler)
}

// Remove all registered handlers from the event.
func (e *Event[Payload]) Flush() {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.handlers = make([]Handler[Payload], 0)
}
