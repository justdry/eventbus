package main

import (
	"context"
	"log"
	"os"

	"github.com/justdry/eventbus"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("you should run in this format:\ngo run . event-name 'json-data'")
	}

	eventName := os.Args[1]
	jsonData := []byte(os.Args[2])

	eventbus.CaptureErrorStack(true)
	bus := setupEventBus()
	bus.Event(eventName).Emit(context.Background(), jsonData)
}

func setupEventBus() *eventbus.EventBus[[]byte] {
	bus := eventbus.New[[]byte]()

	bus.ErrorEvent().Subscribe(func(ctx context.Context, err error, p []byte) {
		sErr, ok := err.(*eventbus.StackedError)
		if !ok {
			log.Fatalf("an error occured: %+v", err)
		}

		log.Fatalf("error occured:\n%s\nstack:\n%s", sErr.Error(), sErr.Stack())
	})

	SetupListeners(bus)
	return bus
}
