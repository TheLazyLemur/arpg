package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/globals"
)

type Bullet struct {
	Position rl.Vector3
	Velocity rl.Vector3
	Lifetime float32
	Speed    float32
	Radius   float32
	Damage   float32
	Active   bool
}

func NewBullet(pos, vel rl.Vector3, speed, lifetime, damage float32) *Bullet {
	return &Bullet{
		Position: pos,
		Velocity: vel,
		Lifetime: lifetime,
		Speed:    speed,
		Radius:   0.1,
		Damage:   damage,
		Active:   true,
	}
}

func (b *Bullet) Update(deltaTime float32) {
	if !b.Active {
		return
	}

	b.Lifetime -= deltaTime
	if b.Lifetime <= 0 {
		b.Active = false
		return
	}

	b.Position.X += b.Velocity.X * b.Speed * deltaTime
	b.Position.Z += b.Velocity.Z * b.Speed * deltaTime
}

func (b *Bullet) IsExpired() bool {
	return !b.Active || b.Lifetime <= 0
}

func (b *Bullet) Deactivate() {
	b.Active = false
}

func (b *Bullet) GetBoundingBox() rl.BoundingBox {
	return rl.BoundingBox{
		Min: rl.Vector3{
			X: b.Position.X - b.Radius,
			Y: b.Position.Y - b.Radius,
			Z: b.Position.Z - b.Radius,
		},
		Max: rl.Vector3{
			X: b.Position.X + b.Radius,
			Y: b.Position.Y + b.Radius,
			Z: b.Position.Z + b.Radius,
		},
	}
}

func (b *Bullet) GetCollisionTags() []string {
	return []string{"bullet"}
}

func (b *Bullet) OnCollision(other globals.Collidable) {
	tags := other.GetCollisionTags()

	for _, tag := range tags {
		switch tag {
		case "enemy":
			entity, ok := other.(*Enemy)
			if !ok {
				panic("Bullet collided with non-enemy entity")
			}
			entity.TakeDamage(b.Damage)

			b.Deactivate()
		case "obstacle":
			b.Deactivate()
		}
	}
}

func (b *Bullet) IsActive() bool {
	return b.Active
}
