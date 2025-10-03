package eventbus

import (
	"sync"
)

func New[Payload any]() EventBus[Payload] {
	return EventBus[Payload]{
		events: make(map[string]*Event[Payload]),
	}
}

type EventBus[Payload any] struct {
	events map[string]*Event[Payload]
	mux    sync.Mutex
}

func (bus *EventBus[Payload]) Event(name string) *Event[Payload] {
	bus.mux.Lock()
	defer bus.mux.Unlock()

	event, exists := bus.events[name]
	if !exists {
		event = newEvent[Payload]()
		bus.events[name] = event
	}

	return event
}
