package eventbus

import (
	"runtime/debug"
	"sync/atomic"
)

var shouldCaptureStack atomic.Bool

// CaptureErrorStack enables or disables stack trace collection
// for newly created StackedError instances.
func CaptureErrorStack(status bool) {
	shouldCaptureStack.Store(status)
}

func NewError(err error) *StackedError {
	if err == nil {
		return nil
	}

	stack := make([]byte, 0)
	if shouldCaptureStack.Load() {
		stack = debug.Stack()
	}

	return &StackedError{
		origin: err,
		stack:  stack,
	}
}

type StackedError struct {
	origin error
	stack  []byte
}

func (e *StackedError) Error() string {
	return e.origin.Error()
}

func (e *StackedError) Unwrap() error {
	return e.origin
}

func (e *StackedError) Stack() []byte {
	return e.stack
}
