# EventBus

<div align="center">

![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue)
[![Go Reference](https://pkg.go.dev/badge/github.com/justdry/eventbus.svg)](https://pkg.go.dev/github.com/justdry/eventbus)
[![Code Coverage](https://codecov.io/github/justdry/eventbus/branch/main/graph/badge.svg?token=C4V0XESSHO)](https://codecov.io/github/justdry/eventbus)
[![Go Quality Report](https://goreportcard.com/badge/github.com/justdry/eventbus)](https://goreportcard.com/report/github.com/justdry/eventbus)
[![Tests](https://github.com/justdry/eventbus/actions/workflows/test.yml/badge.svg)](https://github.com/justdry/eventbus/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Last Release](https://img.shields.io/github/v/release/justdry/eventbus)](https://github.com/justdry/eventbus/releases)

</div>

<p align="center">
	<img src=".github/splash.png" width="60%" alt="JustDRY's EventBus Splash" />
</p>

A small, generic, thread-safe event bus for Go. EventBus is type-safe using Go generics, supports contexts for propagation of request-scoped values and deadlines, and provides a centralized mechanism for handling errors produced by event handlers.

---

## Features

-   Type-safe with Generics — define your own payload type for events.
-   Thread-safe — internal synchronization via mutexes; tests run with the `-race` flag.
-   Centralized Error Handling — all events share a single, attachable error handler (per EventBus instance).
-   Context Support — handlers receive a `context.Context`, so you can carry deadlines, cancellation, or request-scoped values.
-   Parallel Emission — `EmitParallel` triggers all handlers concurrently and waits for all to complete.

---

## Installation

```bash
go get github.com/justdry/eventbus
```

---

## Quick Start

This package exposes a generic EventBus type. Create a bus using `eventbus`.`New[Payload]()`, register named events using `Event(name)` and subscribe handlers to those events.

### Basic example

```go
package main

import (
	"context"
	"fmt"

	"github.com/justdry/eventbus"
)

func main() {
	// Create a new event bus carrying string payloads
	bus := eventbus.New[string]()

	// Get (or create) the named event and register a handler
	greet := bus.Event("greet")
	greet.Subscribe(func(ctx context.Context, name string) error {
		fmt.Printf("Hello %s!", name)
		return nil
	})

	greet.Subscribe(func(ctx context.Context, name string) error {
		fmt.Printf("Hi %s!!!", name)
		return nil
	})

	// Emit the event concurrently and wait until all finish (handlers will run in parallel)
	greet.EmitParallel(context.Background(), "Sina")
}
```

### Centralized error handling

The EventBus creates a shared ErrorEvent for all events. You can subscribe a single error handler to receive errors from handlers across the bus.

```go
bus := eventbus.New[string]()

// Register the handlers that return an error
bus.Event("danger1").Subscribe(func(ctx context.Context, msg string) error {
	// Handlers could wrap returned errors using eventbus.NewError(err) to ensure stack traces are captured.
	return eventbus.NewError(fmt.Errorf("handler failed: %s", msg))
})

bus.Event("danger2").Subscribe(func(ctx context.Context, msg string) error {
	return eventbus.NewError(fmt.Errorf("handler chokhed up: %s", msg))
})

// Register a centralized error handler for the bus
bus.ErrorEvent().Subscribe(func(ctx context.Context, err error, payload string) {
	sErr := err.(*eventbus.StackedError)
	fmt.Printf("Error occurred:\n%s\n%s", sErr.Error(), sErr.Stack())
})

eventbus.CaptureErrorStack(true)

// Emitting will call the handler, and the first returned error (if any)
// will be forwarded to the ErrorEvent handler.
bus.Event("danger1").Emit(context.Background(), "boom!")
bus.Event("danger2").Emit(context.Background(), "boom!")
```

### Using context values in handlers

Handlers receive a context, so you can pass request-scoped data or deadlines.

```go
type ctxKey string
const userKey ctxKey = "user"

bus := eventbus.New[any]()
auth := bus.Event("auth")

auth.Subscribe(func(ctx context.Context, _ any) error {
	user := ctx.Value(userKey).(string)
	fmt.Println("User:", user)
	return nil
})

// Emit with a context that carries the "user"
ctx := context.WithValue(context.Background(), userKey, "sina")
auth.Emit(ctx, nil)
```

---

## Notes

-   EventBus ensures the same ErrorEvent instance is shared across all named events created from a bus, so errors from different events can be handled by a single subscriber.
-   The Event.Emit method iterates over a snapshot of handlers to avoid races with concurrent Subscribe calls; mutexes protect internal state.

---

## Testing

Run the test suite (recommended with the race detector):

```bash
go test ./... -race
```

The tests exercise event emitting, context propagation, error forwarding, and concurrency properties.

---

## Contributing

Contributions, bug reports, and improvements are welcome. Please open issues or PRs on the repository.

## License

MIT License - see [LICENSE](LICENSE) for full text
