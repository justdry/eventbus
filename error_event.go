package eventbus

import (
	"context"
	"sync"
)

func newErrorEvent[Payload any]() *ErrorEvent[Payload] {
	return &ErrorEvent[Payload]{
		handler: nil,
	}
}

type ErrorEvent[Payload any] struct {
	handler ErrorHandler[Payload]
	mux     sync.Mutex
}

// ErrorHandler defines the function signature for handling errors.
type ErrorHandler[Payload any] func(ctx context.Context, err error, p Payload)

// Invoke the subscribed error handler.
// If no handler is subscribed, Emit does nothing.
func (e *ErrorEvent[Payload]) Emit(ctx context.Context, err error, p Payload) {
	e.mux.Lock()
	handler := e.handler
	e.mux.Unlock()

	if handler != nil {
		handler(ctx, err, p)
	}
}

// Register the error handler to be called whenever an error is emitted.
// Any previously registered handler is replaced.
func (e *ErrorEvent[Payload]) Subscribe(handler ErrorHandler[Payload]) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.handler = handler
}

// Remove the subscribed error handler and leaving the handler empty.
func (e *ErrorEvent[Payload]) Flush() {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.handler = nil
}
