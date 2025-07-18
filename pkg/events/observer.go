package events

import (
	"fmt"
	"sync"
)

// Event represents a game event with associated data
type Event struct {
	Type string
	Data interface{}
}

// Observer represents any object that can receive notifications
type Observer interface {
	OnNotify(event Event) error
}

// Subject represents any object that can be observed
type Subject interface {
	Subscribe(eventType string, observer Observer) error
	Unsubscribe(eventType string, observer Observer) error
	Notify(event Event) error
}

// EventBus implements the Subject interface for managing observers
type EventBus struct {
	observers map[string][]Observer
	mutex     sync.RWMutex
}

// NewEventBus creates a new event bus instance
func NewEventBus() *EventBus {
	return &EventBus{
		observers: make(map[string][]Observer),
	}
}

// Subscribe adds an observer for a specific event type
func (eb *EventBus) Subscribe(eventType string, observer Observer) error {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if observer == nil {
		return fmt.Errorf("cannot subscribe nil observer")
	}

	eb.observers[eventType] = append(eb.observers[eventType], observer)
	return nil
}

// Unsubscribe removes an observer from a specific event type
func (eb *EventBus) Unsubscribe(eventType string, observer Observer) error {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	observers, exists := eb.observers[eventType]
	if !exists {
		return nil // Nothing to unsubscribe
	}

	// Find and remove the observer
	for i, obs := range observers {
		if obs == observer {
			eb.observers[eventType] = append(observers[:i], observers[i+1:]...)
			break
		}
	}

	// Clean up empty slices
	if len(eb.observers[eventType]) == 0 {
		delete(eb.observers, eventType)
	}

	return nil
}

// Notify sends an event to all subscribed observers
func (eb *EventBus) Notify(event Event) error {
	eb.mutex.RLock()
	observers, exists := eb.observers[event.Type]
	eb.mutex.RUnlock()

	if !exists {
		return nil // No observers for this event type
	}

	// Create a copy to avoid holding the lock during notification
	observersCopy := make([]Observer, len(observers))
	copy(observersCopy, observers)

	// Notify all observers (outside of lock to prevent deadlocks)
	var errors []error
	for _, observer := range observersCopy {
		if err := observer.OnNotify(event); err != nil {
			errors = append(errors, fmt.Errorf("observer notification failed: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors occurred: %v", errors)
	}

	return nil
}

// GetObserverCount returns the number of observers for a given event type
func (eb *EventBus) GetObserverCount(eventType string) int {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	return len(eb.observers[eventType])
}

// Clear removes all observers
func (eb *EventBus) Clear() {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.observers = make(map[string][]Observer)
}

// GetEventTypes returns all event types that have observers
func (eb *EventBus) GetEventTypes() []string {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	eventTypes := make([]string, 0, len(eb.observers))
	for eventType := range eb.observers {
		eventTypes = append(eventTypes, eventType)
	}
	return eventTypes
}