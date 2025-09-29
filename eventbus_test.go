package eventbus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/justdry/eventbus"

	"github.com/stretchr/testify/assert"
)

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
		name = ctx.Value("name").(string)

		return nil
	})

	ctx := context.WithValue(context.Background(), "name", "sina")
	e.Emit(ctx, "test", nil)

	assert.Equal(t, ctx.Value("name"), name)
}

func TestUseContextInErrorHandler(t *testing.T) {
	var name string

	e := eventbus.NewEventBus[any]()

	e.On("test", func(ctx context.Context, a any) error {
		return errors.New("Test")
	})

	e.OnError(func(ctx context.Context, a any, err error) {
		name = ctx.Value("name").(string)
	})

	ctx := context.WithValue(context.Background(), "name", "sina")
	e.Emit(ctx, "test", nil)
	assert.Equal(t, ctx.Value("name"), name)
}
