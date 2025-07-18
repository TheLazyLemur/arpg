package entities

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/globals"
)

type HealthPickup struct {
	Position   rl.Vector3
	HealAmount float32
	Radius     float32
	Active     bool
}

func NewHealthPickup(pos rl.Vector3) *HealthPickup {
	return &HealthPickup{
		Position:   pos,
		HealAmount: 25.0,
		Radius:     0.3,
		Active:     true,
	}
}

func (h *HealthPickup) Update(deltaTime float32) {}

func (h *HealthPickup) GetBoundingBox() rl.BoundingBox {
	return rl.BoundingBox{
		Min: rl.Vector3{
			X: h.Position.X - h.Radius,
			Y: h.Position.Y - h.Radius,
			Z: h.Position.Z - h.Radius,
		},
		Max: rl.Vector3{
			X: h.Position.X + h.Radius,
			Y: h.Position.Y + h.Radius,
			Z: h.Position.Z + h.Radius,
		},
	}
}

func (h *HealthPickup) GetCollisionTags() []string {
	return []string{"health_pickup"}
}

func (h *HealthPickup) GetTriggerBounds() rl.BoundingBox {
	return h.GetBoundingBox()
}

func (h *HealthPickup) GetTriggerTags() []string {
	return []string{"health_pickup"}
}

func (h *HealthPickup) OnTriggerEnter(other globals.Collidable) {
	if !h.Active {
		return
	}

	tags := other.GetCollisionTags()
	for _, tag := range tags {
		switch tag {
		case "player":
			if player, ok := other.(*Player); ok {
				player.Heal(h.HealAmount)
				h.Active = false
				fmt.Printf("Player healed for %.0f health!\n", h.HealAmount)
			}
		}
	}
}

func (h *HealthPickup) IsActive() bool {
	return h.Active
}
