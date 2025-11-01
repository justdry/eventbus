package eventbus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/justdry/eventbus"
	"github.com/stretchr/testify/require"
)

func TestBusEventAvoidDuplicateEvent(t *testing.T) {
	bus := eventbus.New[any]()

	e1 := bus.Event("test")
	e1.Subscribe(func(ctx context.Context, p any) error { return nil })

	e2 := bus.Event("test")
	e2.Subscribe(func(ctx context.Context, p any) error { return errors.New("") })

	require.Equal(t, e1, e2)
}

// It should run with the `-race` flag
func TestBusEventRaceCondition(t *testing.T) {
	e := eventbus.New[any]()

	registerAndEmit := func() {
		go e.Event("first")

		go e.Event("second")
	}

	require.NotPanics(t, registerAndEmit)
}

func TestBusShouldUseSameErrorEventForAllEvents(t *testing.T) {
	bus := eventbus.New[string]()

	bus.Event("test1").Subscribe(func(ctx context.Context, p string) error {
		return errors.New(" warning in the world")
	})

	bus.Event("test2").Subscribe(func(ctx context.Context, p string) error {
		return errors.New(" error in the world")
	})

	require.Same(t, bus.Event("test1").ErrorEvent, bus.Event("test2").ErrorEvent)
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
	require.Contains(t, trace, "eventbus_test.go")
	require.Contains(t, trace, "justReturnError")
}

func justReturnError(ctx context.Context, p any) error {
	return eventbus.NewError(errors.New("oh!"))
}

func TestDeleteEvent(t *testing.T) {
	bus := eventbus.New[any]()

	t1 := bus.Event("test")
	require.Same(t, t1, bus.Event("test"))

	bus.DeleteEvent("test")
	t2 := bus.Event("test")

	require.NotSame(t, t1, t2)
}
