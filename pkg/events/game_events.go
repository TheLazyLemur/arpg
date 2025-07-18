package events

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game event type constants
const (
	EventTypeBulletSpawn   = "bullet_spawn"
	EventTypeEnemyKilled   = "enemy_killed"
	EventTypePlayerDamaged = "player_damaged"
	EventTypeHealthPickup  = "health_pickup"
	EventTypeGameOver      = "game_over"
	EventTypeVictory       = "victory"
)

// BulletSpawnEvent represents data for bullet spawning
type BulletSpawnEvent struct {
	Position  rl.Vector3
	Direction rl.Vector3
	Speed     float32
	Lifetime  float32
	Damage    float32
}

// EnemyKilledEvent represents data when an enemy is defeated
type EnemyKilledEvent struct {
	EnemyId  string
	Position rl.Vector3
	Killer   string // "player" or "environment"
}

// PlayerDamagedEvent represents data when player takes damage
type PlayerDamagedEvent struct {
	Damage     float32
	Source     string
	NewHealth  float32
	Position   rl.Vector3
}

// HealthPickupEvent represents data when player picks up health
type HealthPickupEvent struct {
	HealAmount  float32
	PickupId    string
	Position    rl.Vector3
	NewHealth   float32
}

// GameOverEvent represents game over state
type GameOverEvent struct {
	Reason      string // "player_died", "time_up", etc.
	FinalScore  int
	Playtime    float32
}

// VictoryEvent represents victory state
type VictoryEvent struct {
	EnemiesKilled int
	TimeElapsed   float32
	Score         int
}

// NewBulletSpawnEvent creates a new bullet spawn event
func NewBulletSpawnEvent(pos, dir rl.Vector3, speed, lifetime, damage float32) Event {
	return Event{
		Type: EventTypeBulletSpawn,
		Data: BulletSpawnEvent{
			Position:  pos,
			Direction: dir,
			Speed:     speed,
			Lifetime:  lifetime,
			Damage:    damage,
		},
	}
}

// NewEnemyKilledEvent creates a new enemy killed event
func NewEnemyKilledEvent(enemyId string, pos rl.Vector3, killer string) Event {
	return Event{
		Type: EventTypeEnemyKilled,
		Data: EnemyKilledEvent{
			EnemyId:  enemyId,
			Position: pos,
			Killer:   killer,
		},
	}
}

// NewPlayerDamagedEvent creates a new player damaged event
func NewPlayerDamagedEvent(damage float32, source string, newHealth float32, pos rl.Vector3) Event {
	return Event{
		Type: EventTypePlayerDamaged,
		Data: PlayerDamagedEvent{
			Damage:    damage,
			Source:    source,
			NewHealth: newHealth,
			Position:  pos,
		},
	}
}

// NewHealthPickupEvent creates a new health pickup event
func NewHealthPickupEvent(healAmount float32, pickupId string, pos rl.Vector3, newHealth float32) Event {
	return Event{
		Type: EventTypeHealthPickup,
		Data: HealthPickupEvent{
			HealAmount: healAmount,
			PickupId:   pickupId,
			Position:   pos,
			NewHealth:  newHealth,
		},
	}
}

// NewGameOverEvent creates a new game over event
func NewGameOverEvent(reason string, finalScore int, playtime float32) Event {
	return Event{
		Type: EventTypeGameOver,
		Data: GameOverEvent{
			Reason:     reason,
			FinalScore: finalScore,
			Playtime:   playtime,
		},
	}
}

// NewVictoryEvent creates a new victory event
func NewVictoryEvent(enemiesKilled int, timeElapsed float32, score int) Event {
	return Event{
		Type: EventTypeVictory,
		Data: VictoryEvent{
			EnemiesKilled: enemiesKilled,
			TimeElapsed:   timeElapsed,
			Score:         score,
		},
	}
}