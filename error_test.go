package eventbus

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	CaptureErrorStack(true)

	origErr := errors.New("original error")
	err := NewError(origErr)

	assert.Contains(t, err.Error(), "original error")
	assert.Contains(t, string(err.Stack()), "error_test.go")
}

func TestErrorUnwrap(t *testing.T) {
	origErr := errors.New("original error")
	err := NewError(origErr)

	assert.ErrorIs(t, err, origErr)

	var target *StackedError
	assert.ErrorAs(t, target, &err)
}

func TestErrorCaptureStack(t *testing.T) {
	origErr := errors.New("original error")

	CaptureErrorStack(false)
	assert.NotContains(t, NewError(origErr).Stack(), "TestErrorCaptureStack")

	CaptureErrorStack(true)
	assert.Contains(t, string(NewError(origErr).Stack()), "TestErrorCaptureStack")
}
