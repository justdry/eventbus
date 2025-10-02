package eventbus

import (
	"context"
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
