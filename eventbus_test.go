package eventbus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/justdry/eventbus"

	"github.com/stretchr/testify/assert"
)

type Example string

const KEY_NAME Example = "name"

func TestEmitEvent(t *testing.T) {
	status := "failed"

	e := eventbus.NewEventBus[any]()

	e.On("test", func(_ context.Context, a any) error {
		status = "passed"

		return nil
	})

	e.Emit(context.Background(), "test", nil)
	assert.NotEqual(t, "failed", status)
}

func TestEmitErrorHandler(t *testing.T) {
	var status string

	e := eventbus.NewEventBus[any]()

	e.On("test", func(_ context.Context, a any) error {
		return errors.New("Test")
	})

	e.OnError(func(_ context.Context, a any, err error) {
		status = "failed"
	})

	err := e.Emit(context.Background(), "test", nil)

	assert.Equal(t, "failed", status)
	assert.NotNil(t, err)
}

func TestEmitUsingContext(t *testing.T) {
	var name string

	e := eventbus.NewEventBus[any]()

	e.On("test", func(ctx context.Context, a any) error {
		name = ctx.Value(KEY_NAME).(string)

		return nil
	})

	ctx := context.WithValue(context.Background(), KEY_NAME, "sina")
	e.Emit(ctx, "test", nil)

	assert.Equal(t, ctx.Value(KEY_NAME), name)
}

func TestUseContextInErrorHandler(t *testing.T) {
	var name string

	e := eventbus.NewEventBus[any]()

	e.On("test", func(ctx context.Context, a any) error {
		return errors.New("Test")
	})

	e.OnError(func(ctx context.Context, a any, err error) {
		name = ctx.Value(KEY_NAME).(string)
	})

	ctx := context.WithValue(context.Background(), KEY_NAME, "sina")
	e.Emit(ctx, "test", nil)
	assert.Equal(t, ctx.Value(KEY_NAME), name)
}

// It should run with the `-race` flag
func TestRaceCondition(t *testing.T) {
	e := eventbus.NewEventBus[any]()

	registerAndEmit := func() {
		go e.On("test", func(ctx context.Context, a any) error {
			return nil
		})

		go e.Emit(context.Background(), "test", nil)
	}

	assert.NotPanics(t, registerAndEmit)
}

func TestEmitAllHandlers(t *testing.T) {
	status := [2]string{"failed", "failed"}

	e := eventbus.NewEventBus[any]()

	e.On("test", func(_ context.Context, a any) error {
		status[0] = "passed 1"

		return nil
	})

	e.On("test", func(_ context.Context, a any) error {
		status[1] = "passed 2"

		return nil
	})

	e.Emit(context.Background(), "test", nil)
	assert.NotEqual(t, "failed", status[0])
	assert.NotEqual(t, "failed", status[1])
}
