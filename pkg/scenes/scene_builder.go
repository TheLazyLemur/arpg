package scenes

import (
	"fmt"

	"arpg/pkg/entities"
	"arpg/pkg/events"
)

// SceneBuilder handles converting JSON scene data to game entities
type SceneBuilder struct{}

// NewSceneBuilder creates a new scene builder
func NewSceneBuilder() *SceneBuilder {
	return &SceneBuilder{}
}

// BuildPlayer creates a player entity from JSON data
func (sb *SceneBuilder) BuildPlayer(data PlayerData, eventBus events.Subject) *entities.Player {
	player := entities.NewPlayer(data.Speed, eventBus)
	player.Position = data.SpawnPoint.ToVector3()
	player.Health = data.Health
	player.MaxHealth = data.MaxHealth
	return player
}

// BuildEnemies creates enemy entities from JSON data
func (sb *SceneBuilder) BuildEnemies(data []EnemyData) []*entities.Enemy {
	enemies := make([]*entities.Enemy, 0, len(data))
	
	for _, enemyData := range data {
		enemy := entities.NewEnemy(
			enemyData.Position.ToVector3(),
			enemyData.Health,
			enemyData.Speed,
		)
		enemies = append(enemies, enemy)
	}
	
	return enemies
}

// BuildObstacles creates obstacle entities from JSON data
func (sb *SceneBuilder) BuildObstacles(data []ObstacleData) []*entities.Obstacle {
	obstacles := make([]*entities.Obstacle, 0, len(data))
	
	for _, obstacleData := range data {
		var obstacle *entities.Obstacle
		color := ParseColor(obstacleData.Color)
		
		switch obstacleData.Type {
		case "box":
			obstacle = entities.NewBoxObstacle(
				obstacleData.Position.ToVector3(),
				obstacleData.Size.ToVector3(),
				color,
			)
		case "cylinder":
			obstacle = entities.NewCylinderObstacle(
				obstacleData.Position.ToVector3(),
				obstacleData.Radius,
				obstacleData.Height,
				color,
			)
		default:
			fmt.Printf("Unknown obstacle type: %s, skipping\\n", obstacleData.Type)
			continue
		}
		
		obstacles = append(obstacles, obstacle)
	}
	
	return obstacles
}

// BuildHealthPickups creates health pickup entities from JSON data
func (sb *SceneBuilder) BuildHealthPickups(data []HealthPickupData) []*entities.HealthPickup {
	pickups := make([]*entities.HealthPickup, 0, len(data))
	
	for _, pickupData := range data {
		pickup := entities.NewHealthPickup(pickupData.Position.ToVector3())
		pickup.HealAmount = pickupData.HealAmount
		pickup.Radius = pickupData.Radius
		pickups = append(pickups, pickup)
	}
	
	return pickups
}

// ValidateSceneData validates the JSON scene data for correctness
func (sb *SceneBuilder) ValidateSceneData(data *SceneData) error {
	// Validate player data
	if data.Player.Speed <= 0 {
		return fmt.Errorf("player speed must be positive, got %f", data.Player.Speed)
	}
	if data.Player.Health <= 0 {
		return fmt.Errorf("player health must be positive, got %f", data.Player.Health)
	}
	if data.Player.MaxHealth <= 0 {
		return fmt.Errorf("player max health must be positive, got %f", data.Player.MaxHealth)
	}
	if data.Player.Health > data.Player.MaxHealth {
		return fmt.Errorf("player health (%f) cannot exceed max health (%f)", 
			data.Player.Health, data.Player.MaxHealth)
	}

	// Validate enemies
	for i, enemy := range data.Entities.Enemies {
		if enemy.ID == "" {
			return fmt.Errorf("enemy %d: ID cannot be empty", i)
		}
		if enemy.Health <= 0 {
			return fmt.Errorf("enemy %s: health must be positive, got %f", enemy.ID, enemy.Health)
		}
		if enemy.Speed <= 0 {
			return fmt.Errorf("enemy %s: speed must be positive, got %f", enemy.ID, enemy.Speed)
		}
	}

	// Validate obstacles
	for i, obstacle := range data.Entities.Obstacles {
		if obstacle.ID == "" {
			return fmt.Errorf("obstacle %d: ID cannot be empty", i)
		}
		if obstacle.Type != "box" && obstacle.Type != "cylinder" {
			return fmt.Errorf("obstacle %s: type must be 'box' or 'cylinder', got %s", 
				obstacle.ID, obstacle.Type)
		}
		if obstacle.Type == "box" {
			if obstacle.Size.X <= 0 || obstacle.Size.Y <= 0 || obstacle.Size.Z <= 0 {
				return fmt.Errorf("obstacle %s: box size must be positive, got (%f, %f, %f)", 
					obstacle.ID, obstacle.Size.X, obstacle.Size.Y, obstacle.Size.Z)
			}
		}
		if obstacle.Type == "cylinder" {
			if obstacle.Radius <= 0 {
				return fmt.Errorf("obstacle %s: cylinder radius must be positive, got %f", 
					obstacle.ID, obstacle.Radius)
			}
			if obstacle.Height <= 0 {
				return fmt.Errorf("obstacle %s: cylinder height must be positive, got %f", 
					obstacle.ID, obstacle.Height)
			}
		}
	}

	// Validate health pickups
	for i, pickup := range data.Entities.HealthPickups {
		if pickup.ID == "" {
			return fmt.Errorf("health pickup %d: ID cannot be empty", i)
		}
		if pickup.HealAmount <= 0 {
			return fmt.Errorf("health pickup %s: heal amount must be positive, got %f", 
				pickup.ID, pickup.HealAmount)
		}
		if pickup.Radius <= 0 {
			return fmt.Errorf("health pickup %s: radius must be positive, got %f", 
				pickup.ID, pickup.Radius)
		}
	}

	return nil
}

// CheckIDUniqueness ensures all entity IDs are unique within the scene
func (sb *SceneBuilder) CheckIDUniqueness(data *SceneData) error {
	ids := make(map[string]bool)
	
	// Check enemy IDs
	for _, enemy := range data.Entities.Enemies {
		if ids[enemy.ID] {
			return fmt.Errorf("duplicate ID found: %s", enemy.ID)
		}
		ids[enemy.ID] = true
	}
	
	// Check obstacle IDs
	for _, obstacle := range data.Entities.Obstacles {
		if ids[obstacle.ID] {
			return fmt.Errorf("duplicate ID found: %s", obstacle.ID)
		}
		ids[obstacle.ID] = true
	}
	
	// Check health pickup IDs
	for _, pickup := range data.Entities.HealthPickups {
		if ids[pickup.ID] {
			return fmt.Errorf("duplicate ID found: %s", pickup.ID)
		}
		ids[pickup.ID] = true
	}
	
	return nil
}