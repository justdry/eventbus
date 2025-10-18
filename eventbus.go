package eventbus

import (
	"sync"
)

// Create a new EventBus instance.
func New[Payload any]() *EventBus[Payload] {
	return &EventBus[Payload]{
		events:     make(map[string]*Event[Payload]),
		errorEvent: newErrorEvent[Payload](),
	}
}

type EventBus[Payload any] struct {
	events     map[string]*Event[Payload]
	errorEvent *ErrorEvent[Payload]
	mux        sync.Mutex
}

// Return the Event with the given name.
//
// If the event does not already exist, it will be created.
func (bus *EventBus[Payload]) Event(name string) *Event[Payload] {
	bus.mux.Lock()
	defer bus.mux.Unlock()

	event, exists := bus.events[name]
	if !exists {
		event = newEvent[Payload]()
		event.ErrorEvent = bus.errorEvent

		bus.events[name] = event
	}

	return event
}

// Return the shared ErrorEvent for the EventBus.
//
// This event is emitted whenever a handler returns an error.
func (bus *EventBus[Payload]) ErrorEvent() *ErrorEvent[Payload] {
	return bus.errorEvent
}
