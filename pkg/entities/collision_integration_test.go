package entities

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/events"
)

func TestEnemyPlayerCollisionDamageEvents(t *testing.T) {
	// given
	// ... a mock event bus
	// ... a player with the event bus
	// ... an enemy positioned to collide with the player
	mockEventBus := &MockEventBus{}
	player := NewPlayer(5.0, mockEventBus)
	player.Position = rl.Vector3{X: 0, Y: 0, Z: 0}
	player.Health = 100.0
	enemy := NewEnemy(rl.Vector3{X: 0, Y: 0, Z: 0}, 50.0, 2.0)
	enemy.AttackCooldown = 0 // Ready to attack
	// when
	// ... the enemy collides with the player
	enemy.OnCollision(player)
	// then
	// ... should emit exactly one player damage event
	// ... and should have correct damage values
	// ... and player should have reduced health
	// ... and enemy should have attack cooldown set
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
	if damageData.Damage != 10.0 {
		t.Fatalf("Expected damage 10.0, got %f", damageData.Damage)
	}
	if damageData.NewHealth != 90.0 {
		t.Fatalf("Expected new health 90.0, got %f", damageData.NewHealth)
	}
	if player.Health != 90.0 {
		t.Fatalf("Expected player health 90.0, got %f", player.Health)
	}
	if enemy.AttackCooldown != 1.0 {
		t.Fatalf("Expected attack cooldown 1.0, got %f", enemy.AttackCooldown)
	}
}

func TestEnemyPlayerCollisionCooldownPreventsSpam(t *testing.T) {
	// given
	// ... a mock event bus
	// ... a player with the event bus
	// ... an enemy with attack cooldown active
	mockEventBus := &MockEventBus{}
	player := NewPlayer(5.0, mockEventBus)
	player.Position = rl.Vector3{X: 0, Y: 0, Z: 0}
	player.Health = 100.0
	enemy := NewEnemy(rl.Vector3{X: 0, Y: 0, Z: 0}, 50.0, 2.0)
	enemy.AttackCooldown = 0.5 // Still on cooldown
	// when
	// ... the enemy collides with the player while on cooldown
	enemy.OnCollision(player)
	// then
	// ... should not emit any events
	// ... and player should not take damage
	if len(mockEventBus.events) != 0 {
		t.Fatalf("Expected 0 events due to cooldown, got %d", len(mockEventBus.events))
	}
	if player.Health != 100.0 {
		t.Fatalf("Expected player health unchanged at 100.0, got %f", player.Health)
	}
}

func TestEnemyPlayerFatalCollisionEmitsGameOver(t *testing.T) {
	// given
	// ... a mock event bus
	// ... a player with very low health
	// ... an enemy ready to attack
	mockEventBus := &MockEventBus{}
	player := NewPlayer(5.0, mockEventBus)
	player.Position = rl.Vector3{X: 0, Y: 0, Z: 0}
	player.Health = 5.0 // Low health - will die from 10 damage
	enemy := NewEnemy(rl.Vector3{X: 0, Y: 0, Z: 0}, 50.0, 2.0)
	enemy.AttackCooldown = 0
	// when
	// ... the enemy deals fatal damage to the player
	enemy.OnCollision(player)
	// then
	// ... should emit two events: damage and game over
	// ... first event should be player damage
	// ... second event should be game over
	// ... and player should be dead
	if len(mockEventBus.events) != 2 {
		t.Fatalf("Expected 2 events (damage + game over), got %d", len(mockEventBus.events))
	}
	damageEvent := mockEventBus.events[0]
	if damageEvent.Type != events.EventTypePlayerDamaged {
		t.Fatal("First event should be player damage")
	}
	gameOverEvent := mockEventBus.events[1]
	if gameOverEvent.Type != events.EventTypeGameOver {
		t.Fatal("Second event should be game over")
	}
	if player.IsAlive() {
		t.Fatal("Player should be dead")
	}
	if player.Health != 0.0 {
		t.Fatalf("Expected player health 0.0, got %f", player.Health)
	}
}

