package eventbus

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorShouldImplementGoInterface(t *testing.T) {
	getError := func() {
		var err error = newError[any](errors.New("an error"), nil)

		assert.Equal(t, "an error", err.Error())
	}

	assert.NotPanics(t, getError)
}

func TestErrorHasCorrectPayload(t *testing.T) {
	err := newError[int32](errors.New("an error"), 1)

	assert.Equal(t, int32(1), err.Payload)
}
