package entities

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/globals"
)

type Enemy struct {
	Position       rl.Vector3
	Health         float32
	MaxHealth      float32
	Radius         float32
	Height         float32
	Speed          float32
	Active         bool
	Target         *Player
	AttackCooldown float32
}

func NewEnemy(pos rl.Vector3, health, speed float32) *Enemy {
	return &Enemy{
		Position:  pos,
		Health:    health,
		MaxHealth: health,
		Radius:    0.6,
		Height:    1.5,
		Speed:     speed,
		Active:    true,
	}
}

func (e *Enemy) Update(deltaTime float32, player *Player) {
	if !e.Active || player == nil {
		return
	}

	if e.AttackCooldown > 0 {
		e.AttackCooldown -= deltaTime
	}

	deltaX := player.Position.X - e.Position.X
	deltaZ := player.Position.Z - e.Position.Z

	distance := float32(math.Sqrt(float64(deltaX*deltaX + deltaZ*deltaZ)))

	if distance > 1.0 {
		dirX := deltaX / distance
		dirZ := deltaZ / distance

		e.Position.X += dirX * e.Speed * deltaTime
		e.Position.Z += dirZ * e.Speed * deltaTime
	}
}

func (e *Enemy) GetBoundingBox() rl.BoundingBox {
	return rl.BoundingBox{
		Min: rl.Vector3{
			X: e.Position.X - e.Radius,
			Y: e.Position.Y,
			Z: e.Position.Z - e.Radius,
		},
		Max: rl.Vector3{
			X: e.Position.X + e.Radius,
			Y: e.Position.Y + e.Height,
			Z: e.Position.Z + e.Radius,
		},
	}
}

func (e *Enemy) TakeDamage(damage float32) {
	e.Health -= damage
	if e.Health <= 0 {
		e.Health = 0
		e.Active = false
	}
}

func (e *Enemy) IsAlive() bool {
	return e.Active && e.Health > 0
}

func (e *Enemy) GetHealthPercent() float32 {
	if e.MaxHealth == 0 {
		return 0
	}
	return e.Health / e.MaxHealth
}

func (e *Enemy) GetHeadPosition() rl.Vector3 {
	return rl.Vector3{
		X: e.Position.X,
		Y: e.Position.Y + 1.0,
		Z: e.Position.Z,
	}
}

func (e *Enemy) GetHealthBarPosition() rl.Vector3 {
	return rl.Vector3{
		X: e.Position.X,
		Y: e.Position.Y + 2.0,
		Z: e.Position.Z,
	}
}

func (e *Enemy) GetCollisionTags() []string {
	return []string{"enemy"}
}

func (e *Enemy) OnCollision(other globals.Collidable) {
	tags := other.GetCollisionTags()
	for _, tag := range tags {
		switch tag {
		case "player":
			if e.AttackCooldown <= 0 {
				p := other.(*Player)
				p.TakeDamage(10.0)
				e.AttackCooldown = 1.0
			}
		}
	}
}

func (e *Enemy) IsActive() bool {
	return e.IsAlive()
}
