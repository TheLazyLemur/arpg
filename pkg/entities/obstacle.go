package entities

import (
	"fmt"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/globals"
)

type ObstacleType int

const (
	ObstacleTypeBox ObstacleType = iota
	ObstacleTypeCylinder
)

type Obstacle struct {
	Position rl.Vector3
	Size     rl.Vector3
	Radius   float32
	Height   float32
	Type     ObstacleType
	Color    rl.Color
	Active   bool
}

func NewBoxObstacle(pos rl.Vector3, size rl.Vector3, color rl.Color) *Obstacle {
	return &Obstacle{
		Position: pos,
		Size:     size,
		Type:     ObstacleTypeBox,
		Color:    color,
		Active:   true,
	}
}

func NewCylinderObstacle(pos rl.Vector3, radius, height float32, color rl.Color) *Obstacle {
	return &Obstacle{
		Position: pos,
		Radius:   radius,
		Height:   height,
		Type:     ObstacleTypeCylinder,
		Color:    color,
		Active:   true,
	}
}

func (o *Obstacle) GetBoundingBox() rl.BoundingBox {
	switch o.Type {
	case ObstacleTypeBox:
		return rl.BoundingBox{
			Min: rl.Vector3{
				X: o.Position.X - o.Size.X/2,
				Y: o.Position.Y - o.Size.Y/2,
				Z: o.Position.Z - o.Size.Z/2,
			},
			Max: rl.Vector3{
				X: o.Position.X + o.Size.X/2,
				Y: o.Position.Y + o.Size.Y/2,
				Z: o.Position.Z + o.Size.Z/2,
			},
		}
	case ObstacleTypeCylinder:
		return rl.BoundingBox{
			Min: rl.Vector3{
				X: o.Position.X - o.Radius,
				Y: o.Position.Y,
				Z: o.Position.Z - o.Radius,
			},
			Max: rl.Vector3{
				X: o.Position.X + o.Radius,
				Y: o.Position.Y + o.Height,
				Z: o.Position.Z + o.Radius,
			},
		}
	default:
		return rl.BoundingBox{}
	}
}

func (o *Obstacle) IsBox() bool {
	return o.Type == ObstacleTypeBox
}

func (o *Obstacle) IsCylinder() bool {
	return o.Type == ObstacleTypeCylinder
}

func (o *Obstacle) GetCollisionTags() []string {
	return []string{"obstacle"}
}

func (o *Obstacle) OnCollision(other globals.Collidable) {
	if slices.Contains(other.GetCollisionTags(), "player") {
		o.Active = false
		fmt.Println("Collision with obstacle detected:", other.GetCollisionTags())
	}
}

func (o *Obstacle) IsActive() bool {
	return o.Active // Obstacles are always active
}
