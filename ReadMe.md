# Event Bus

<img src=".github/splash.jpg" width="60%" alt="JustDRY's EventBus Splash" />

A simple, generic, thread-safe event bus for Go.
`eventbus` lets you register event handlers, emit events, and hook into error handling using Go’s `context.Context`.

---

## ✨ Features

-   **Type-safe with Generics** — define your own event type.
-   **Thread-safe** — tested with Go’s `-race` flag.
-   **Error Handling** — attach an error handler to centralize error management.
-   **Context Support** — propagate deadlines, cancelation signals, or request-scoped data into your handlers.

---

## 🚀 Installation

```bash
go get github.com/justdry/eventbus
```

---

## 📖 Usage

### Basic Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/justdry/eventbus"
)

func main() {
	// Create a new event bus
	e := eventbus.NewEventBus[string]()

	// Register a handler
	e.On("greet", func(ctx context.Context, name string) error {
		fmt.Printf("Hello, %s!\n", name)
		return nil
	})

	// Emit an event
	e.Emit(context.Background(), "greet", "Sina")
}
```

---

### Error Handling

```go
e := eventbus.NewEventBus[string]()

// Register a handler that may fail
e.On("fail", func(ctx context.Context, msg string) error {
	return fmt.Errorf("something went wrong: %s", msg)
})

// Register a centralized error handler
e.OnError(func(ctx context.Context, msg string, err error) {
	fmt.Println("Error occurred:", err)
})

// Emit the event
_ = e.Emit(context.Background(), "fail", "boom!")
```

---

### Using Context

```go
type ctxKey string

const userKey ctxKey = "user"

e := eventbus.NewEventBus[any]()

// Register handler using context value
e.On("auth", func(ctx context.Context, _ any) error {
	user := ctx.Value(userKey).(string)
	fmt.Println("User:", user)
	return nil
})

// Emit with context
ctx := context.WithValue(context.Background(), userKey, "sina")
e.Emit(ctx, "auth", nil)
```

---

## ✅ Testing

This package includes a full test suite. Run:

```bash
go test ./... -race
```

The `-race` flag is recommended to ensure thread safety.
