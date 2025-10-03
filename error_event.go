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

type ErrorHandler[Payload any] func(ctx context.Context, err error, p Payload)

func (e *ErrorEvent[Payload]) Emit(ctx context.Context, err error, p Payload) {
	e.mux.Lock()
	handler := e.handler
	e.mux.Unlock()

	if handler != nil {
		handler(ctx, err, p)
	}
}

func (e *ErrorEvent[Payload]) Subscribe(handler ErrorHandler[Payload]) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.handler = handler
}

func (e *ErrorEvent[Payload]) Flush() {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.handler = nil
}
