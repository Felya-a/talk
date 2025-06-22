package core

import "sync"

type Event interface {
	Name() string
}

type EventHandler func(e Event)

// EventBus реализует простую событийную шину
type EventBus struct {
	subscribers map[string][]EventHandler
	lock        sync.RWMutex
}

// NewEventBus создаёт новую шину событий
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe подписывает на событие с именем eventName
func (eb *EventBus) Subscribe(eventName string, handler EventHandler) {
	eb.lock.Lock()
	defer eb.lock.Unlock()
	eb.subscribers[eventName] = append(eb.subscribers[eventName], handler)
}

// Publish публикует событие всем подписчикам
func (eb *EventBus) Publish(event Event) {
	eb.lock.RLock()
	defer eb.lock.RUnlock()
	for _, handler := range eb.subscribers[event.Name()] {
		go handler(event)
	}
}
