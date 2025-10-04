package eventbus

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Example string

const KEY_NAME Example = "name"

func TestEmitEvent(t *testing.T) {
	status := "failed"

	e := newEvent[any]()

	e.Subscribe(func(_ context.Context, a any) error {
		status = "passed"

		return nil
	})

	e.Emit(context.Background(), nil)
	assert.NotEqual(t, "failed", status)
}

func TestEmitUsingContext(t *testing.T) {
	var name string

	e := newEvent[any]()

	e.Subscribe(func(ctx context.Context, a any) error {
		name = ctx.Value(KEY_NAME).(string)

		return nil
	})

	ctx := context.WithValue(context.Background(), KEY_NAME, "sina")
	e.Emit(ctx, nil)

	assert.Equal(t, ctx.Value(KEY_NAME), name)
}

// It should run with the `-race` flag
func TestEventRaceCondition(t *testing.T) {
	e := newEvent[any]()

	registerAndEmit := func() {
		go e.Subscribe(func(ctx context.Context, a any) error {
			return nil
		})

		go e.Emit(context.Background(), nil)
	}

	assert.NotPanics(t, registerAndEmit)
}

func TestEmitAllHandlers(t *testing.T) {
	status := [2]string{"failed", "failed"}

	e := newEvent[any]()

	e.Subscribe(func(_ context.Context, a any) error {
		status[0] = "passed 1"

		return nil
	})

	e.Subscribe(func(_ context.Context, a any) error {
		status[1] = "passed 2"

		return nil
	})

	e.Emit(context.Background(), nil)
	assert.NotEqual(t, "failed", status[0])
	assert.NotEqual(t, "failed", status[1])
}

func TestEventEmitsErrorHandlerOnReceivingError(t *testing.T) {
	status := "bad"

	event := newEvent[string]()
	event.Subscribe(func(ctx context.Context, p string) error {
		return errors.New("Oh")
	})

	errorEvent := newErrorEvent[string]()
	errorEvent.Subscribe(func(ctx context.Context, err error, p string) {
		status = p
	})

	assert.NotPanics(t, func() {
		event.Emit(context.Background(), "world")
	})

	event.ErrorEvent = errorEvent

	err := event.Emit(context.Background(), "world")
	assert.NotEqual(t, "Oh", err)

	assert.Equal(t, "world", status)
}

func TestFlushEventHandlers(t *testing.T) {
	status := "nothing changed"

	event := newEvent[string]()
	event.Subscribe(func(ctx context.Context, p string) error {
		status = p
		return nil
	})

	event.Flush()

	event.Emit(context.Background(), "Hello Universe!")
	assert.Equal(t, "nothing changed", status)
}
