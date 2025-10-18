package eventbus

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	CaptureErrorStack(true)

	origErr := errors.New("original error")
	err := NewError(origErr)

	require.Contains(t, err.Error(), "original error")
	require.Contains(t, string(err.Stack()), "error_test.go")
}

func TestErrorUnwrap(t *testing.T) {
	origErr := errors.New("original error")
	err := NewError(origErr)

	require.ErrorIs(t, err, origErr)

	var target *StackedError
	require.ErrorAs(t, target, &err)
}

func TestErrorCaptureStack(t *testing.T) {
	origErr := errors.New("original error")

	CaptureErrorStack(false)
	require.NotContains(t, NewError(origErr).Stack(), "TestErrorCaptureStack")

	CaptureErrorStack(true)
	require.Contains(t, string(NewError(origErr).Stack()), "TestErrorCaptureStack")
}
