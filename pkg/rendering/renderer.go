package rendering

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/camera"
	"arpg/pkg/config"
	"arpg/pkg/entities"
)

type Renderer struct {
	config *config.Config
}

func NewRenderer(cfg *config.Config) *Renderer {
	return &Renderer{
		config: cfg,
	}
}

func (r *Renderer) Initialize() {
	rl.InitWindow(r.config.Window.Width, r.config.Window.Height, r.config.Window.Title)

	if r.config.Window.VSync {
		rl.SetTargetFPS(r.config.Window.TargetFPS)
	}

	if r.config.Window.Fullscreen {
		rl.ToggleFullscreen()
	}
}

func (r *Renderer) Shutdown() {
	rl.CloseWindow()
}

func (r *Renderer) BeginFrame() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
}

func (r *Renderer) EndFrame() {
	rl.EndDrawing()
}

func (r *Renderer) BeginMode3D(cam *camera.Camera) {
	rl.BeginMode3D(cam.GetRaylibCamera())
}

func (r *Renderer) EndMode3D() {
	rl.EndMode3D()
}

func (r *Renderer) DrawGround() {
	groundSize := float32(50.0)
	rl.DrawPlane(
		rl.Vector3{X: 0, Y: -0.5, Z: 0},
		rl.Vector2{X: groundSize, Y: groundSize},
		rl.White,
	)
}

func (r *Renderer) DrawGrid() {
	if r.config.Graphics.DrawGrid {
		rl.DrawGrid(20, 1.0)
	}
}

func (r *Renderer) DrawHealthPickups(healthPickups []*entities.HealthPickup) {
	for _, pickup := range healthPickups {
		if pickup.IsActive() {
			rl.DrawSphere(pickup.Position, 0.3, rl.Green)
		}
	}
}

func (r *Renderer) DrawPlayer(player *entities.Player) {
	rl.DrawCylinder(player.Position, player.Radius, player.Radius, player.Height, 8, rl.Blue)

	headPos := rl.Vector3{
		X: player.Position.X,
		Y: player.Position.Y + 3,
		Z: player.Position.Z,
	}
	rl.DrawSphere(headPos, 0.2, rl.Red)

	r.drawGun(player)
}

func (r *Renderer) drawGun(player *entities.Player) {
	gunLength := float32(0.8)
	gunRadius := float32(0.05)

	gunStart := rl.Vector3{
		X: player.Position.X,
		Y: player.Position.Y + 1.0,
		Z: player.Position.Z,
	}

	gunEnd := rl.Vector3{
		X: player.Position.X + gunLength*float32(math.Cos(float64(player.Rotation))),
		Y: player.Position.Y + 1.0,
		Z: player.Position.Z + gunLength*float32(math.Sin(float64(player.Rotation))),
	}

	rl.DrawLine3D(gunStart, gunEnd, rl.Black)
	rl.DrawSphere(gunEnd, gunRadius, rl.DarkGray)
}

func (r *Renderer) DrawBullets(bullets []*entities.Bullet) {
	for _, bullet := range bullets {
		if bullet.Active {
			rl.DrawSphere(bullet.Position, bullet.Radius, rl.Yellow)
		}
	}
}

func (r *Renderer) DrawEnemies(enemies []*entities.Enemy) {
	for _, enemy := range enemies {
		if enemy.IsAlive() {
			rl.DrawCylinder(enemy.Position, enemy.Radius, enemy.Radius, enemy.Height, 8, rl.Red)

			headPos := enemy.GetHeadPosition()
			rl.DrawSphere(headPos, 0.3, rl.Black)
		}
	}
}

func (r *Renderer) DrawObstacles(obstacles []*entities.Obstacle) {
	for _, obstacle := range obstacles {
		if !obstacle.Active {
			continue
		}
		switch obstacle.Type {
		case entities.ObstacleTypeBox:
			rl.DrawCube(
				obstacle.Position,
				obstacle.Size.X,
				obstacle.Size.Y,
				obstacle.Size.Z,
				obstacle.Color,
			)
			if r.config.Graphics.DrawWires {
				rl.DrawCubeWires(
					obstacle.Position,
					obstacle.Size.X,
					obstacle.Size.Y,
					obstacle.Size.Z,
					rl.DarkBrown,
				)
			}
		case entities.ObstacleTypeCylinder:
			rl.DrawCylinder(
				obstacle.Position,
				obstacle.Radius,
				obstacle.Radius,
				obstacle.Height,
				8,
				obstacle.Color,
			)
			if r.config.Graphics.DrawWires {
				rl.DrawCylinderWires(
					obstacle.Position,
					obstacle.Radius,
					obstacle.Radius,
					obstacle.Height,
					8,
					rl.Green,
				)
			}
		}
	}
}

func (r *Renderer) DrawEnemyHealthBars(enemies []*entities.Enemy, cam *camera.Camera) {
	if !r.config.Debug.ShowHealthBars {
		return
	}

	for _, enemy := range enemies {
		if !enemy.IsAlive() {
			continue
		}

		healthBarPos := enemy.GetHealthBarPosition()
		screenPos := rl.GetWorldToScreen(healthBarPos, cam.GetRaylibCamera())

		if screenPos.X >= 0 && screenPos.X <= float32(r.config.Window.Width) &&
			screenPos.Y >= 0 && screenPos.Y <= float32(r.config.Window.Height) {

			barWidth := float32(50)
			barHeight := float32(6)

			rl.DrawRectangle(
				int32(screenPos.X-barWidth/2),
				int32(screenPos.Y),
				int32(barWidth),
				int32(barHeight),
				rl.Red,
			)

			healthPercent := enemy.GetHealthPercent()
			rl.DrawRectangle(
				int32(screenPos.X-barWidth/2),
				int32(screenPos.Y),
				int32(barWidth*healthPercent),
				int32(barHeight),
				rl.Green,
			)

			rl.DrawRectangleLines(
				int32(screenPos.X-barWidth/2),
				int32(screenPos.Y),
				int32(barWidth),
				int32(barHeight),
				rl.Black,
			)
		}
	}
}

func (r *Renderer) DrawUI() {
	rl.DrawText("WASD to move", 10, 10, 20, rl.DarkGray)
	rl.DrawText("Mouse to aim", 10, 35, 20, rl.DarkGray)
	rl.DrawText("Left click to shoot", 10, 60, 20, rl.DarkGray)

	if r.config.Debug.ShowFPS {
		rl.DrawFPS(10, 90)
	}
}

func (r *Renderer) ShouldClose() bool {
	return rl.WindowShouldClose()
}
