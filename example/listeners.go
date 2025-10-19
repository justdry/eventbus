package main

import (
	"context"
	"fmt"

	"github.com/justdry/eventbus"
	"github.com/justdry/eventbus/example/handlers"
)

func SetupListeners(bus *eventbus.EventBus[[]byte]) {
	greeting := bus.Event("greeting")
	greeting.Subscribe(handlers.Hi)
	greeting.Subscribe(handlers.Greeting)

	bus.Event("weather").Subscribe(handlers.Weather)

	bus.Event("error").Subscribe(func(ctx context.Context, p []byte) error {
		return fmt.Errorf("there is no any stack trace")
	})
}
