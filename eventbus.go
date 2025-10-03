package eventbus

import (
	"sync"
)

func New[Payload any]() EventBus[Payload] {
	return EventBus[Payload]{
		events:     make(map[string]*Event[Payload]),
		errorEvent: newErrorEvent[Payload](),
	}
}

type EventBus[Payload any] struct {
	events     map[string]*Event[Payload]
	errorEvent *ErrorEvent[Payload]
	mux        sync.Mutex
}

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

func (bus *EventBus[Payload]) ErrorEvent() *ErrorEvent[Payload] {
	return bus.errorEvent
}
