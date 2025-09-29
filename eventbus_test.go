package eventbus_test

import (
	"errors"
	"testing"

	"github.com/justdry/eventbus"

	"github.com/stretchr/testify/assert"
)

func TestEmitEvent(t *testing.T) {
	status := "failed"

	e := eventbus.NewEventBus[any]()

	e.On("test", func(a any) error {
		status = "passed"

		return nil
	})

	e.Emit("test", nil)
	assert.NotEqual(t, "failed", status)
}

func TestEmitErrorHandler(t *testing.T) {
	var status string

	e := eventbus.NewEventBus[any]()

	e.On("test", func(a any) error {
		return errors.New("Test")
	})

	e.OnError(func(event any, err error) {
		status = "failed"
	})

	err := e.Emit("test", nil)

	assert.Equal(t, "failed", status)
	assert.NotNil(t, err)
}
