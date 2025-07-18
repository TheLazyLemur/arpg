package globals

import "sync"

type EventHandler func(data any)

type Event interface {
	RegisterHandler(eventType string, handler EventHandler)
	EmitEvent(eventType string, data any)
}

var EventSystem Event

type DefaultEvent struct {
	handlers map[string][]EventHandler
	lock     sync.Mutex
}

func InitEvent() {
	EventSystem = &DefaultEvent{}
}

func (de *DefaultEvent) RegisterHandler(eventType string, handler EventHandler) {
	de.lock.Lock()
	defer de.lock.Unlock()

	if de.handlers == nil {
		de.handlers = make(map[string][]EventHandler)
	}

	de.handlers[eventType] = append(de.handlers[eventType], handler)
}

func (de *DefaultEvent) EmitEvent(eventType string, data any) {
	de.lock.Lock()
	defer de.lock.Unlock()

	if de.handlers == nil {
		de.handlers = make(map[string][]EventHandler)
	}

	handlers, exists := de.handlers[eventType]
	if exists {
		for _, handler := range handlers {
			handler(data)
		}
	}
}
