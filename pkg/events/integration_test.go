package events

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TestObserverPatternIntegration demonstrates the complete Observer pattern workflow
func TestObserverPatternIntegration(t *testing.T) {
	// given
	// ... a new event bus
	// ... a mock game scene that implements Observer
	// ... a mock player that publishes events
	eventBus := NewEventBus()
	scene := &MockGameScene{}

	// when
	// ... the scene subscribes to bullet spawn events
	err := eventBus.Subscribe(EventTypeBulletSpawn, scene)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// ... and a bullet spawn event is created and notified
	bulletEvent := NewBulletSpawnEvent(
		rl.Vector3{X: 1, Y: 0, Z: 1},
		rl.Vector3{X: 0, Y: 0, Z: 1},
		10.0,
		5.0,
		25.0,
	)
	err = eventBus.Notify(bulletEvent)
	if err != nil {
		t.Fatalf("Failed to notify: %v", err)
	}

	// then
	// ... the scene should have received the bullet spawn event
	if len(scene.receivedEvents) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(scene.receivedEvents))
	}

	receivedEvent := scene.receivedEvents[0]
	if receivedEvent.Type != EventTypeBulletSpawn {
		t.Fatalf("Expected event type %s, got %s", EventTypeBulletSpawn, receivedEvent.Type)
	}

	bulletData, ok := receivedEvent.Data.(BulletSpawnEvent)
	if !ok {
		t.Fatal("Event data should be BulletSpawnEvent")
	}

	if bulletData.Speed != 10.0 {
		t.Fatalf("Expected speed 10.0, got %f", bulletData.Speed)
	}
}

// TestMultipleObserversIntegration tests multiple observers receiving the same event
func TestMultipleObserversIntegration(t *testing.T) {
	// given
	// ... a new event bus
	// ... multiple mock observers
	eventBus := NewEventBus()
	scene1 := &MockGameScene{}
	scene2 := &MockGameScene{}
	ui := &MockUIObserver{}
	// when
	// ... all observers subscribe to player damage events
	// ... and a player damage event is notified
	eventBus.Subscribe(EventTypePlayerDamaged, scene1)
	eventBus.Subscribe(EventTypePlayerDamaged, scene2)
	eventBus.Subscribe(EventTypePlayerDamaged, ui)
	damageEvent := NewPlayerDamagedEvent(
		25.0,
		"enemy",
		75.0,
		rl.Vector3{X: 0, Y: 0, Z: 0},
	)
	err := eventBus.Notify(damageEvent)
	if err != nil {
		t.Fatalf("Failed to notify: %v", err)
	}
	// then
	// ... all observers should receive the event
	// ... and all should receive the same event data
	if len(scene1.receivedEvents) != 1 {
		t.Fatal("Scene1 should receive 1 event")
	}
	if len(scene2.receivedEvents) != 1 {
		t.Fatal("Scene2 should receive 1 event")
	}
	if len(ui.receivedEvents) != 1 {
		t.Fatal("UI should receive 1 event")
	}
	for _, observer := range []*MockGameScene{scene1, scene2} {
		event := observer.receivedEvents[0]
		if event.Type != EventTypePlayerDamaged {
			t.Fatal("Should receive player damaged event")
		}

		damageData, ok := event.Data.(PlayerDamagedEvent)
		if !ok || damageData.Damage != 25.0 {
			t.Fatal("Should receive correct damage data")
		}
	}
}

// MockGameScene implements Observer for testing
type MockGameScene struct {
	receivedEvents []Event
}

func (m *MockGameScene) OnNotify(event Event) error {
	m.receivedEvents = append(m.receivedEvents, event)
	return nil
}

// MockUIObserver implements Observer for testing UI components
type MockUIObserver struct {
	receivedEvents []Event
}

func (m *MockUIObserver) OnNotify(event Event) error {
	m.receivedEvents = append(m.receivedEvents, event)
	return nil
}

