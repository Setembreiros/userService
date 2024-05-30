package bus

import (
	"context"
)

type Event struct {
	Type string
	Data []byte
}

type EventBus struct {
	subscribers map[string][]chan<- Event
}

type EventSubscription struct {
	EventType string
	Handler   EventHandler
}

type EventHandler interface {
	Handle(event []byte)
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan<- Event),
	}
}

func (eb *EventBus) Publish(event Event) {
	subscriberChannels := eb.subscribers[event.Type]

	for _, subscriberChannel := range subscriberChannels {
		subscriberChannel <- event
	}
}

func (eb *EventBus) Subscribe(subscription *EventSubscription, ctx context.Context) {
	subscriptionChan := make(chan Event)
	eb.subscribers[subscription.EventType] = append(eb.subscribers[subscription.EventType], subscriptionChan)
	go subscription.handle(subscriptionChan, ctx)
}

func (es EventSubscription) handle(busChannel <-chan Event, ctx context.Context) {
	for {
		select {
		case event := <-busChannel:
			go es.Handler.Handle(event.Data)
		case <-ctx.Done():
			return
		}
	}
}
