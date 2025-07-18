package entities

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/events"
)

// MockEventBus for testing player events
type MockEventBus struct {
	events []events.Event
}

func (m *MockEventBus) Subscribe(eventType string, observer events.Observer) error {
	return nil
}

func (m *MockEventBus) Unsubscribe(eventType string, observer events.Observer) error {
	return nil
}

func (m *MockEventBus) Notify(event events.Event) error {
	m.events = append(m.events, event)
	return nil
}

func TestPlayerTakeDamageEmitsEvents(t *testing.T) {
	// given
	// ... a player with a mock event bus
	// ... the player starts with full health
	mockEventBus := &MockEventBus{}
	player := NewPlayer(5.0, mockEventBus)
	player.Health = 100.0
	player.MaxHealth = 100.0
	// when
	// ... the player takes damage but doesn't die
	player.TakeDamage(25.0)
	// then
	// ... should emit exactly one player damaged event
	// ... and should be a player damaged event
	// ... and should have correct damage data
	if len(mockEventBus.events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(mockEventBus.events))
	}
	event := mockEventBus.events[0]
	if event.Type != events.EventTypePlayerDamaged {
		t.Fatalf("Expected player damaged event, got %s", event.Type)
	}
	damageData, ok := event.Data.(events.PlayerDamagedEvent)
	if !ok {
		t.Fatal("Event data should be PlayerDamagedEvent")
	}
	if damageData.Damage != 25.0 {
		t.Fatalf("Expected damage 25.0, got %f", damageData.Damage)
	}
	if damageData.NewHealth != 75.0 {
		t.Fatalf("Expected new health 75.0, got %f", damageData.NewHealth)
	}
	if damageData.Source != "enemy" {
		t.Fatalf("Expected source 'enemy', got %s", damageData.Source)
	}
}

func TestPlayerDeathEmitsGameOverEvent(t *testing.T) {
	// given
	// ... a player with a mock event bus
	// ... the player starts with low health
	mockEventBus := &MockEventBus{}
	player := NewPlayer(5.0, mockEventBus)
	player.Health = 10.0
	player.MaxHealth = 100.0
	// when
	// ... the player takes fatal damage
	player.TakeDamage(15.0)
	// then
	// ... should emit two events: player damaged and game over
	// ... first event should be player damaged
	// ... second event should be game over
	// ... and should have correct game over data
	// ... and player should be dead
	if len(mockEventBus.events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(mockEventBus.events))
	}
	damageEvent := mockEventBus.events[0]
	if damageEvent.Type != events.EventTypePlayerDamaged {
		t.Fatal("First event should be player damaged")
	}
	gameOverEvent := mockEventBus.events[1]
	if gameOverEvent.Type != events.EventTypeGameOver {
		t.Fatal("Second event should be game over")
	}
	gameOverData, ok := gameOverEvent.Data.(events.GameOverEvent)
	if !ok {
		t.Fatal("Event data should be GameOverEvent")
	}
	if gameOverData.Reason != "player_died" {
		t.Fatalf("Expected reason 'player_died', got %s", gameOverData.Reason)
	}
	if player.Health != 0.0 {
		t.Fatalf("Expected health 0.0, got %f", player.Health)
	}
	if player.IsAlive() {
		t.Fatal("Player should not be alive")
	}
}

// MockCamera for testing
type MockCamera struct {
	worldPos rl.Vector3
}

func (m *MockCamera) GetWorldPositionFromMouse(mousePos rl.Vector2) rl.Vector3 {
	return m.worldPos
}
