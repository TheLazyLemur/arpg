package camera

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/config"
	"arpg/pkg/entities"
)

type Camera struct {
	camera rl.Camera3D
	offset rl.Vector3
	config *config.Config
}

func NewCamera(cfg *config.Config) *Camera {
	return &Camera{
		offset: rl.Vector3{X: 0, Y: 10, Z: 8},
		config: cfg,
	}
}

func (c *Camera) Initialize(player *entities.Player) {
	c.camera = rl.Camera3D{
		Position: rl.Vector3{
			X: player.Position.X + c.offset.X,
			Y: player.Position.Y + c.offset.Y,
			Z: player.Position.Z + c.offset.Z,
		},
		Target:     player.Position,               // Look at player
		Up:         rl.Vector3{X: 0, Y: 0, Z: -1}, // Changed Up vector for top-down
		Fovy:       c.config.Graphics.FOV,
		Projection: rl.CameraPerspective,
	}
}

func (c *Camera) Update(player *entities.Player) {
	// Update camera position to follow the player
	c.camera.Position = rl.Vector3{
		X: player.Position.X + c.offset.X,
		Y: player.Position.Y + c.offset.Y,
		Z: player.Position.Z + c.offset.Z,
	}
	c.camera.Target = player.Position
}

func (c *Camera) GetRaylibCamera() rl.Camera3D {
	return c.camera
}

func (c *Camera) GetWorldPositionFromMouse(mousePos rl.Vector2) rl.Vector3 {
	// Cast a ray from the camera through the mouse position
	ray := rl.GetMouseRay(mousePos, c.camera)

	// Find intersection with ground plane (Y = 0)
	if ray.Direction.Y != 0 {
		t := -ray.Position.Y / ray.Direction.Y
		return rl.Vector3{
			X: ray.Position.X + t*ray.Direction.X,
			Y: 0,
			Z: ray.Position.Z + t*ray.Direction.Z,
		}
	}

	return rl.Vector3{X: 0, Y: 0, Z: 0}
}

func (c *Camera) SetOffset(offset rl.Vector3) {
	c.offset = offset
}

func (c *Camera) GetOffset() rl.Vector3 {
	return c.offset
}

func (c *Camera) SetFOV(fov float32) {
	c.camera.Fovy = fov
	c.config.Graphics.FOV = fov
}

