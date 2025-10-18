package eventbus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/justdry/eventbus"
	"github.com/stretchr/testify/assert"
)

func TestBusEventAvoidDuplicateEvent(t *testing.T) {
	bus := eventbus.New[any]()

	e1 := bus.Event("test")
	e1.Subscribe(func(ctx context.Context, p any) error { return nil })

	e2 := bus.Event("test")
	e2.Subscribe(func(ctx context.Context, p any) error { return errors.New("") })

	assert.Equal(t, e1, e2)
}

// It should run with the `-race` flag
func TestBusEventRaceCondition(t *testing.T) {
	e := eventbus.New[any]()

	registerAndEmit := func() {
		go e.Event("first")

		go e.Event("second")
	}

	assert.NotPanics(t, registerAndEmit)
}

func TestBusShouldUseSameErrorEventForAllEvents(t *testing.T) {
	bus := eventbus.New[string]()

	bus.Event("test1").Subscribe(func(ctx context.Context, p string) error {
		return errors.New(" warning in the world")
	})

	bus.Event("test2").Subscribe(func(ctx context.Context, p string) error {
		return errors.New(" error in the world")
	})

	assert.Same(t, bus.Event("test1").ErrorEvent, bus.Event("test2").ErrorEvent)
}

func TestErrorHandlerStackTrace(t *testing.T) {
	trace := ""

	bus := eventbus.New[any]()

	bus.ErrorEvent().Subscribe(func(ctx context.Context, err error, p any) {
		trace = string(err.(*eventbus.StackedError).Stack())
	})

	test := bus.Event("test")

	test.Subscribe(justReturnError)

	eventbus.CaptureErrorStack(true)

	test.Emit(context.Background(), nil)
	assert.Contains(t, trace, "eventbus_test.go")
	assert.Contains(t, trace, "justReturnError")
}

func justReturnError(ctx context.Context, p any) error {
	return eventbus.NewError(errors.New("oh!"))
}
