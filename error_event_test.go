package eventbus

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmitErrorEvent(t *testing.T) {
	status := "Don't know"

	e := newErrorEvent[any]()

	e.Subscribe(func(_ context.Context, err error, a any) {
		status = err.Error()
	})

	e.Emit(context.Background(), errors.New("I know"), nil)
	require.Equal(t, "I know", status)
}

func TestEmitErrorEventUsingContext(t *testing.T) {
	var name string

	e := newErrorEvent[any]()

	e.Subscribe(func(ctx context.Context, err error, a any) {
		name = ctx.Value(KEY_NAME).(string)
	})

	ctx := context.WithValue(context.Background(), KEY_NAME, "sina")
	e.Emit(ctx, nil, nil)

	require.Equal(t, ctx.Value(KEY_NAME), name)
}

// It should run with the `-race` flag
func TestErrorEventRaceCondition(t *testing.T) {
	e := newErrorEvent[any]()

	registerAndEmit := func() {
		go e.Subscribe(func(ctx context.Context, err error, a any) {
		})

		go e.Emit(context.Background(), nil, nil)
	}

	require.NotPanics(t, registerAndEmit)
}

func TestFlushErrorHandler(t *testing.T) {
	status := "Don't know"

	e := newErrorEvent[any]()

	e.Subscribe(func(_ context.Context, err error, a any) {
		status = err.Error()
	})

	e.Flush()

	e.Emit(context.Background(), errors.New("I know"), nil)
	require.Equal(t, "Don't know", status)
}
