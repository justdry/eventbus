package eventbus_test

import (
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
