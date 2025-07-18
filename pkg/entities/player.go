package entities

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/events"
	"arpg/pkg/globals"
)

type Player struct {
	Position  rl.Vector3
	Rotation  float32
	Speed     float32
	Radius    float32
	Height    float32
	Health    float32
	MaxHealth float32
	eventBus  events.Subject // Injected dependency for events
}

func NewPlayer(speed float32, eventBus events.Subject) *Player {
	return &Player{
		Position:  rl.Vector3{X: 0, Y: 0, Z: 0},
		Rotation:  0,
		Speed:     speed,
		Radius:    0.5,
		Height:    1.0,
		Health:    100.0,
		MaxHealth: 100.0,
		eventBus:  eventBus,
	}
}

type CameraInterface interface {
	GetWorldPositionFromMouse(mousePos rl.Vector2) rl.Vector3
}

func (p *Player) Update(
	deltaTime float32,
	obstacles []*Obstacle,
	camera CameraInterface,
) {
	p.updateRotation(camera)

	p.updateMovement(deltaTime, obstacles)

	if globals.InputSystem.IsMouseLeftPressed() {
		p.shoot(camera)
	}
}

func (p *Player) shoot(camera CameraInterface) {
	if p.eventBus == nil {
		return // Cannot shoot without event bus
	}

	gunTip := p.GetGunTip()

	mousePos := globals.InputSystem.GetMousePosition()
	worldPos := camera.GetWorldPositionFromMouse(mousePos)

	direction := rl.Vector3{
		X: worldPos.X - gunTip.X,
		Y: 0,
		Z: worldPos.Z - gunTip.Z,
	}

	direction = rl.Vector3Normalize(direction)

	// Create and emit bullet spawn event using the new Observer pattern
	bulletEvent := events.NewBulletSpawnEvent(
		gunTip,
		direction,
		15.0, // Speed
		3.0,  // Lifetime
		25.0, // Damage
	)

	// Notify observers of the bullet spawn event
	if err := p.eventBus.Notify(bulletEvent); err != nil {
		// Log error but don't prevent gameplay
		fmt.Printf("Error notifying bullet spawn: %v\n", err)
	}
}

func (p *Player) updateRotation(camera CameraInterface) {
	mousePos := globals.InputSystem.GetMousePosition()
	worldPos := camera.GetWorldPositionFromMouse(mousePos)

	deltaX := worldPos.X - p.Position.X
	deltaZ := worldPos.Z - p.Position.Z
	p.Rotation = float32(math.Atan2(float64(deltaZ), float64(deltaX)))
}

func (p *Player) updateMovement(deltaTime float32, obstacles []*Obstacle) {
	moveSpeed := p.Speed * deltaTime

	if globals.InputSystem.IsUpDown() {
		newPos := p.Position
		newPos.Z -= moveSpeed
		if globals.Collision.CheckMovement(p, newPos) {
			p.Position = newPos
		}
	}

	if globals.InputSystem.IsDownDown() {
		newPos := p.Position
		newPos.Z += moveSpeed
		if globals.Collision.CheckMovement(p, newPos) {
			p.Position = newPos
		}
	}

	if globals.InputSystem.IsLeftDown() {
		newPos := p.Position
		newPos.X -= moveSpeed
		if globals.Collision.CheckMovement(p, newPos) {
			p.Position = newPos
		}
	}

	if globals.InputSystem.IsRightDown() {
		newPos := p.Position
		newPos.X += moveSpeed
		if globals.Collision.CheckMovement(p, newPos) {
			p.Position = newPos
		}
	}
}

func (p *Player) GetBoundingBox() rl.BoundingBox {
	return rl.BoundingBox{
		Min: rl.Vector3{
			X: p.Position.X - p.Radius,
			Y: p.Position.Y,
			Z: p.Position.Z - p.Radius,
		},
		Max: rl.Vector3{
			X: p.Position.X + p.Radius,
			Y: p.Position.Y + p.Height,
			Z: p.Position.Z + p.Radius,
		},
	}
}

func (p *Player) GetCollisionTags() []string {
	return []string{"player"}
}

func (p *Player) OnCollision(other globals.Collidable) {}

func (p *Player) IsActive() bool {
	return p.IsAlive()
}

func (p *Player) GetGunTip() rl.Vector3 {
	gunLength := float32(0.8)
	return rl.Vector3{
		X: p.Position.X + gunLength*float32(math.Cos(float64(p.Rotation))),
		Y: p.Position.Y + 1.0,
		Z: p.Position.Z + gunLength*float32(math.Sin(float64(p.Rotation))),
	}
}

func (p *Player) GetDirection() rl.Vector3 {
	return rl.Vector3{
		X: float32(math.Cos(float64(p.Rotation))),
		Y: 0,
		Z: float32(math.Sin(float64(p.Rotation))),
	}
}

func (p *Player) TakeDamage(damage float32) {
	oldHealth := p.Health
	p.Health -= damage
	if p.Health < 0 {
		p.Health = 0
	}

	// Emit player damaged event if event bus is available
	if p.eventBus != nil {
		damageEvent := events.NewPlayerDamagedEvent(
			damage,
			"enemy", // TODO: Make this dynamic based on damage source
			p.Health,
			p.Position,
		)

		if err := p.eventBus.Notify(damageEvent); err != nil {
			fmt.Printf("Error notifying player damage: %v\n", err)
		}

		// Emit game over event if player died
		if oldHealth > 0 && p.Health <= 0 {
			gameOverEvent := events.NewGameOverEvent(
				"player_died",
				0, // TODO: Implement score system
				0, // TODO: Implement playtime tracking
			)

			if err := p.eventBus.Notify(gameOverEvent); err != nil {
				fmt.Printf("Error notifying game over: %v\n", err)
			}
		}
	}
}

func (p *Player) IsAlive() bool {
	return p.Health > 0
}

func (p *Player) Heal(amount float32) {
	p.Health += amount
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
}
