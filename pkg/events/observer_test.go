package events

import (
	"errors"
	"sync"
	"testing"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Mock observer for testing
type MockObserver struct {
	receivedEvents []Event
	shouldError    bool
	mutex          sync.Mutex
}

func (m *MockObserver) OnNotify(event Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.receivedEvents = append(m.receivedEvents, event)

	if m.shouldError {
		return errors.New("mock observer error")
	}

	return nil
}

func (m *MockObserver) GetReceivedEvents() []Event {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	events := make([]Event, len(m.receivedEvents))
	copy(events, m.receivedEvents)
	return events
}

func (m *MockObserver) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.receivedEvents = nil
}

func TestEventBusCreation(t *testing.T) {
	// given
	// ... no pre-existing event bus
	// when
	// ... a new event bus is created
	eventBus := NewEventBus()
	// then
	// ... should have a non-nil event bus
	// ... and should have empty observers map
	// ... and should have zero event types
	if eventBus == nil {
		t.Fatal("NewEventBus should return non-nil event bus")
	}
	if len(eventBus.observers) != 0 {
		t.Fatal("New event bus should have empty observers map")
	}
	if len(eventBus.GetEventTypes()) != 0 {
		t.Fatal("New event bus should have no event types")
	}
}

func TestEventBusSubscription(t *testing.T) {
	// given
	// ... a new event bus
	// ... a mock observer
	eventBus := NewEventBus()
	observer := &MockObserver{}
	// when
	// ... the observer subscribes to an event type
	err := eventBus.Subscribe("test_event", observer)
	// then
	// ... should complete without error
	// ... and should have one observer for the event type
	// ... and should have the event type in the list
	if err != nil {
		t.Fatalf("Subscribe should not return error: %v", err)
	}
	if eventBus.GetObserverCount("test_event") != 1 {
		t.Fatalf("Expected 1 observer, got %d", eventBus.GetObserverCount("test_event"))
	}
	eventTypes := eventBus.GetEventTypes()
	if len(eventTypes) != 1 || eventTypes[0] != "test_event" {
		t.Fatal("Event type should be in the list")
	}
	// when
	// ... a second observer subscribes to the same event type
	observer2 := &MockObserver{}
	err = eventBus.Subscribe("test_event", observer2)
	// then
	// ... should complete without error
	// ... and should have two observers for the event type
	if err != nil {
		t.Fatalf("Second subscribe should not return error: %v", err)
	}
	if eventBus.GetObserverCount("test_event") != 2 {
		t.Fatalf("Expected 2 observers, got %d", eventBus.GetObserverCount("test_event"))
	}
}

func TestEventBusSubscriptionNilObserver(t *testing.T) {
	// given
	// ... a new event bus
	// ... a nil observer
	eventBus := NewEventBus()
	// when
	// ... attempting to subscribe a nil observer
	err := eventBus.Subscribe("test_event", nil)
	// then
	// ... should return an error
	// ... and should have zero observers
	if err == nil {
		t.Fatal("Subscribe with nil observer should return error")
	}
	if eventBus.GetObserverCount("test_event") != 0 {
		t.Fatal("Should have zero observers after nil subscription")
	}
}

func TestEventBusUnsubscription(t *testing.T) {
	// given
	// ... a new event bus
	// ... two observers subscribed to the same event type
	eventBus := NewEventBus()
	observer1 := &MockObserver{}
	observer2 := &MockObserver{}
	eventBus.Subscribe("test_event", observer1)
	eventBus.Subscribe("test_event", observer2)
	// when
	// ... the first observer unsubscribes
	err := eventBus.Unsubscribe("test_event", observer1)
	// then
	// ... should complete without error
	// ... and should have one observer remaining
	if err != nil {
		t.Fatalf("Unsubscribe should not return error: %v", err)
	}
	if eventBus.GetObserverCount("test_event") != 1 {
		t.Fatalf("Expected 1 observer remaining, got %d", eventBus.GetObserverCount("test_event"))
	}
	// when
	// ... the second observer unsubscribes
	err = eventBus.Unsubscribe("test_event", observer2)
	// then
	// ... should complete without error
	// ... and should have zero observers
	// ... and should have no event types
	if err != nil {
		t.Fatalf("Second unsubscribe should not return error: %v", err)
	}
	if eventBus.GetObserverCount("test_event") != 0 {
		t.Fatalf("Expected 0 observers, got %d", eventBus.GetObserverCount("test_event"))
	}
	if len(eventBus.GetEventTypes()) != 0 {
		t.Fatal("Should have no event types after all unsubscriptions")
	}
}

func TestEventBusUnsubscriptionNonExistent(t *testing.T) {
	// given
	// ... a new event bus
	// ... an observer that was never subscribed
	eventBus := NewEventBus()
	observer := &MockObserver{}
	// when
	// ... attempting to unsubscribe a non-existent observer
	err := eventBus.Unsubscribe("test_event", observer)
	// then
	// ... should complete without error
	// ... and should have zero observers
	if err != nil {
		t.Fatalf("Unsubscribe non-existent should not return error: %v", err)
	}
	if eventBus.GetObserverCount("test_event") != 0 {
		t.Fatal("Should have zero observers")
	}
}

func TestEventBusNotification(t *testing.T) {
	// given
	// ... a new event bus
	// ... two observers subscribed to the same event type
	// ... a test event
	eventBus := NewEventBus()
	observer1 := &MockObserver{}
	observer2 := &MockObserver{}
	eventBus.Subscribe("test_event", observer1)
	eventBus.Subscribe("test_event", observer2)
	testEvent := Event{Type: "test_event", Data: "test_data"}
	// when
	// ... the event is notified
	err := eventBus.Notify(testEvent)
	// then
	// ... should complete without error
	// ... and both observers should receive the event
	// ... and should receive the correct event
	if err != nil {
		t.Fatalf("Notify should not return error: %v", err)
	}
	events1 := observer1.GetReceivedEvents()
	events2 := observer2.GetReceivedEvents()
	if len(events1) != 1 {
		t.Fatalf("Observer1 should receive 1 event, got %d", len(events1))
	}
	if len(events2) != 1 {
		t.Fatalf("Observer2 should receive 1 event, got %d", len(events2))
	}
	if events1[0].Type != "test_event" || events1[0].Data != "test_data" {
		t.Fatal("Observer1 received incorrect event")
	}
	if events2[0].Type != "test_event" || events2[0].Data != "test_data" {
		t.Fatal("Observer2 received incorrect event")
	}
}

func TestEventBusNotificationNoObservers(t *testing.T) {
	// given
	// ... a new event bus with no observers
	// ... a test event
	eventBus := NewEventBus()
	testEvent := Event{Type: "test_event", Data: "test_data"}
	// when
	// ... the event is notified
	err := eventBus.Notify(testEvent)
	// then
	// ... should complete without error
	if err != nil {
		t.Fatalf("Notify with no observers should not return error: %v", err)
	}
}

func TestEventBusNotificationWithErrors(t *testing.T) {
	// given
	// ... a new event bus
	// ... an observer that returns errors
	// ... a normal observer
	eventBus := NewEventBus()
	errorObserver := &MockObserver{shouldError: true}
	normalObserver := &MockObserver{}
	eventBus.Subscribe("test_event", errorObserver)
	eventBus.Subscribe("test_event", normalObserver)
	testEvent := Event{Type: "test_event", Data: "test_data"}
	// when
	// ... the event is notified
	err := eventBus.Notify(testEvent)
	// then
	// ... should return an error
	// ... but normal observer should still receive the event
	if err == nil {
		t.Fatal("Notify should return error when observer fails")
	}
	events := normalObserver.GetReceivedEvents()
	if len(events) != 1 {
		t.Fatal("Normal observer should still receive event despite other observer errors")
	}
}

func TestEventBusClear(t *testing.T) {
	// given
	// ... a new event bus
	// ... multiple observers subscribed to different event types
	eventBus := NewEventBus()
	observer1 := &MockObserver{}
	observer2 := &MockObserver{}
	eventBus.Subscribe("event1", observer1)
	eventBus.Subscribe("event2", observer2)
	// when
	// ... the event bus is cleared
	eventBus.Clear()
	// then
	// ... should have no observers for any event type
	// ... and should have no event types
	if eventBus.GetObserverCount("event1") != 0 {
		t.Fatal("Should have zero observers for event1 after clear")
	}
	if eventBus.GetObserverCount("event2") != 0 {
		t.Fatal("Should have zero observers for event2 after clear")
	}
	if len(eventBus.GetEventTypes()) != 0 {
		t.Fatal("Should have no event types after clear")
	}
}

func TestEventBusConcurrency(t *testing.T) {
	// given
	// ... a new event bus
	// ... multiple goroutines subscribing, unsubscribing, and notifying
	eventBus := NewEventBus()
	numGoroutines := 100
	eventsPerGoroutine := 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // Subscribe, notify, unsubscribe
	// when
	// ... multiple goroutines perform concurrent operations
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			observer := &MockObserver{}
			eventBus.Subscribe("concurrent_test", observer)
		}(i)

		go func(id int) {
			defer wg.Done()
			for j := range eventsPerGoroutine {
				event := Event{Type: "concurrent_test", Data: id*eventsPerGoroutine + j}
				eventBus.Notify(event)
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Millisecond) // Let subscribe happen first
			observer := &MockObserver{}
			eventBus.Unsubscribe("concurrent_test", observer)
		}(i)
	}
	// then
	// ... should complete without deadlocks or panics
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		// Success - no deadlocks
	case <-time.After(5 * time.Second):
		t.Fatal("Concurrent operations timed out - possible deadlock")
	}
}

func TestGameEventCreation(t *testing.T) {
	// given
	// ... bullet spawn parameters
	position := rl.Vector3{X: 1, Y: 2, Z: 3}
	direction := rl.Vector3{X: 0, Y: 0, Z: 1}
	speed := float32(10.0)
	lifetime := float32(5.0)
	damage := float32(25.0)
	// when
	// ... a bullet spawn event is created
	event := NewBulletSpawnEvent(position, direction, speed, lifetime, damage)
	// then
	// ... should have correct event type
	// ... and should have correct data
	if event.Type != EventTypeBulletSpawn {
		t.Fatalf("Expected event type %s, got %s", EventTypeBulletSpawn, event.Type)
	}
	bulletData, ok := event.Data.(BulletSpawnEvent)
	if !ok {
		t.Fatal("Event data should be BulletSpawnEvent")
	}
	if bulletData.Position != position {
		t.Fatal("Position should match")
	}
	if bulletData.Direction != direction {
		t.Fatal("Direction should match")
	}
	if bulletData.Speed != speed {
		t.Fatal("Speed should match")
	}
	if bulletData.Lifetime != lifetime {
		t.Fatal("Lifetime should match")
	}
	if bulletData.Damage != damage {
		t.Fatal("Damage should match")
	}
}

