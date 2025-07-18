package scenes

import (
	"fmt"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"

	"arpg/pkg/camera"
	"arpg/pkg/config"
	"arpg/pkg/entities"
	"arpg/pkg/events"
	"arpg/pkg/globals"
	"arpg/pkg/rendering"
	"arpg/pkg/scenes"
)

type WorldScene struct {
	config *config.Config

	// Systems
	camera       *camera.Camera
	sceneBuilder *scenes.SceneBuilder
	eventBus     *events.EventBus

	// Game entities
	player        *entities.Player
	enemies       []*entities.Enemy
	bullets       []*entities.Bullet
	obstacles     []*entities.Obstacle
	healthPickups []*entities.HealthPickup

	// Scene state
	shouldTransition bool
	nextScene        string
	paused           bool
}

func NewGameScene(cfg *config.Config) *WorldScene {
	gs := &WorldScene{
		config:           cfg,
		shouldTransition: false,
		nextScene:        "",
		paused:           false,
		sceneBuilder:     scenes.NewSceneBuilder(),
		eventBus:         events.NewEventBus(),
	}

	globals.InitInput()
	globals.InitCollision()
	globals.InitTriggers()
	gs.camera = camera.NewCamera(cfg)

	return gs
}

func (gs *WorldScene) Initialize() error {
	// Clear any existing subscriptions to prevent duplicate event handling
	gs.eventBus.Clear()
	
	if err := gs.eventBus.Subscribe(events.EventTypeBulletSpawn, gs); err != nil {
		return fmt.Errorf("failed to subscribe to bullet spawn events: %w", err)
	}

	if err := gs.eventBus.Subscribe(events.EventTypePlayerDamaged, gs); err != nil {
		return fmt.Errorf("failed to subscribe to player damaged events: %w", err)
	}

	if err := gs.eventBus.Subscribe(events.EventTypeGameOver, gs); err != nil {
		return fmt.Errorf("failed to subscribe to game over events: %w", err)
	}

	return gs.InitializeFromJSON("scenes/game_scene.json")
}

func (gs *WorldScene) OnNotify(event events.Event) error {
	switch event.Type {
	case events.EventTypeBulletSpawn:
		bulletData, ok := event.Data.(events.BulletSpawnEvent)
		if !ok {
			return fmt.Errorf("invalid bullet spawn event data")
		}
		gs.spawnBulletFromEvent(bulletData)

	case events.EventTypePlayerDamaged:
		damageData, ok := event.Data.(events.PlayerDamagedEvent)
		if ok {
			log.Printf("Player took %.1f damage from %s, health: %.1f",
				damageData.Damage, damageData.Source, damageData.NewHealth)
		}

	case events.EventTypeGameOver:
		gameOverData, ok := event.Data.(events.GameOverEvent)
		if ok {
			log.Printf("Game over: %s", gameOverData.Reason)
		}
	}

	return nil
}

func (gs *WorldScene) InitializeFromJSON(jsonFile string) error {
	globals.Collision.ClearAll()
	globals.Triggers.ClearAll()

	sceneData, err := scenes.LoadSceneFromJSON(jsonFile)
	if err != nil {
		return err
	}

	if err := gs.sceneBuilder.ValidateSceneData(sceneData); err != nil {
		return err
	}

	if err := gs.sceneBuilder.CheckIDUniqueness(sceneData); err != nil {
		return err
	}

	gs.player = gs.sceneBuilder.BuildPlayer(sceneData.Player, gs.eventBus)
	gs.enemies = gs.sceneBuilder.BuildEnemies(sceneData.Entities.Enemies)
	gs.obstacles = gs.sceneBuilder.BuildObstacles(sceneData.Entities.Obstacles)
	gs.healthPickups = gs.sceneBuilder.BuildHealthPickups(sceneData.Entities.HealthPickups)
	gs.bullets = make([]*entities.Bullet, 0)

	globals.Collision.RegisterCollidable(gs.player)
	globals.Triggers.RegisterCollidable(gs.player) // Player can activate triggers
	for _, enemy := range gs.enemies {
		globals.Collision.RegisterCollidable(enemy)
	}
	for _, obstacle := range gs.obstacles {
		globals.Collision.RegisterCollidable(obstacle)
	}
	for _, pickup := range gs.healthPickups {
		globals.Triggers.RegisterTrigger(pickup)
	}

	gs.camera.Initialize(gs.player)

	gs.shouldTransition = false
	gs.nextScene = ""
	gs.paused = false

	log.Printf("Scene loaded successfully: %s", sceneData.Metadata.Name)
	return nil
}

func (gs *WorldScene) Update(deltaTime float32) error {
	if gs.paused {
		return nil
	}

	gs.updateEntities(deltaTime)

	globals.Collision.Update()
	globals.Triggers.Update()

	gs.cleanupEntities()
	gs.camera.Update(gs.player)

	return nil
}

func (gs *WorldScene) Render(renderer *rendering.Renderer) error {
	renderer.BeginFrame()

	renderer.BeginMode3D(gs.camera)

	renderer.DrawGround()
	renderer.DrawObstacles(gs.obstacles)
	renderer.DrawPlayer(gs.player)
	renderer.DrawEnemies(gs.enemies)
	renderer.DrawBullets(gs.bullets)
	renderer.DrawHealthPickups(gs.healthPickups)
	renderer.DrawGrid()

	renderer.EndMode3D()

	renderer.DrawEnemyHealthBars(gs.enemies, gs.camera)
	gs.drawGameUI()

	renderer.EndFrame()
	return nil
}

func (gs *WorldScene) drawGameUI() {
	rl.DrawText("WASD to move", 10, 10, 20, rl.DarkGray)
	rl.DrawText("Mouse to aim", 10, 35, 20, rl.DarkGray)
	rl.DrawText("Left click to shoot", 10, 60, 20, rl.DarkGray)
	rl.DrawText("ESC to return to menu", 10, 85, 20, rl.DarkGray)

	if gs.config.Debug.ShowFPS {
		rl.DrawFPS(10, 110)
	}

	healthText := "Health: " + floatToString(gs.player.Health, 0)
	rl.DrawText(healthText, 10, gs.config.Window.Height-60, 20, rl.Red)

	aliveEnemies := 0
	for _, enemy := range gs.enemies {
		if enemy.IsAlive() {
			aliveEnemies++
		}
	}
	enemyText := "Enemies: " + intToString(aliveEnemies)
	rl.DrawText(enemyText, 10, gs.config.Window.Height-35, 20, rl.Blue)

	activeBullets := 0
	for _, bullet := range gs.bullets {
		if bullet.Active {
			activeBullets++
		}
	}
	bulletText := "Bullets: " + intToString(activeBullets)
	rl.DrawText(bulletText, 10, gs.config.Window.Height-10, 20, rl.Yellow)

	if gs.paused {
		pauseText := "PAUSED - Press P to resume"
		pauseFontSize := int32(32)
		pauseWidth := rl.MeasureText(pauseText, pauseFontSize)
		pauseX := (gs.config.Window.Width - pauseWidth) / 2
		pauseY := gs.config.Window.Height / 2

		rl.DrawText(pauseText, pauseX, pauseY, pauseFontSize, rl.White)
	}

	if !gs.player.IsAlive() {
		gameOverText := "GAME OVER - Press R to restart or ESC for menu"
		gameOverFontSize := int32(24)
		gameOverWidth := rl.MeasureText(gameOverText, gameOverFontSize)
		gameOverX := (gs.config.Window.Width - gameOverWidth) / 2
		gameOverY := gs.config.Window.Height / 2

		rl.DrawText(gameOverText, gameOverX, gameOverY, gameOverFontSize, rl.Red)
	}

	if aliveEnemies == 0 {
		victoryText := "VICTORY! - Press R to restart or ESC for menu"
		victoryFontSize := int32(24)
		victoryWidth := rl.MeasureText(victoryText, victoryFontSize)
		victoryX := (gs.config.Window.Width - victoryWidth) / 2
		victoryY := gs.config.Window.Height / 2

		rl.DrawText(victoryText, victoryX, victoryY, victoryFontSize, rl.Green)
	}
}

func (gs *WorldScene) HandleInput(deltaTime float32) error {
	// Handle pause
	if globals.InputSystem.IsPausePressed() {
		gs.paused = !gs.paused
	}

	if globals.InputSystem.IsSpacePressed() {
		gs.shouldTransition = true
		gs.nextScene = "menu"
		return nil
	}

	if globals.InputSystem.IsRestartPressed() {
		return gs.Initialize() // Restart the game
	}

	if globals.InputSystem.IsDebugPressed() {
		gs.config.Debug.ShowFPS = !gs.config.Debug.ShowFPS
	}

	return nil
}

func (gs *WorldScene) Cleanup() error {
	gs.enemies = nil
	gs.bullets = nil
	gs.obstacles = nil
	gs.player = nil

	return nil
}

func (gs *WorldScene) GetName() string {
	return "game"
}

func (gs *WorldScene) ShouldTransition() bool {
	return gs.shouldTransition
}

func (gs *WorldScene) GetNextScene() string {
	return gs.nextScene
}

func (gs *WorldScene) spawnBulletFromEvent(data events.BulletSpawnEvent) {
	bullet := entities.NewBullet(
		data.Position,
		data.Direction,
		data.Speed,
		data.Lifetime,
		data.Damage,
	)

	gs.bullets = append(gs.bullets, bullet)

	globals.Collision.RegisterCollidable(bullet)
}

func (gs *WorldScene) updateEntities(deltaTime float32) {
	for _, enemy := range gs.enemies {
		enemy.Update(deltaTime, gs.player)
	}

	for _, bullet := range gs.bullets {
		bullet.Update(deltaTime)
	}

	for _, pickup := range gs.healthPickups {
		pickup.Update(deltaTime)
	}

	if gs.player.IsAlive() {
		gs.player.Update(deltaTime, gs.obstacles, gs.camera)
	}
}

func (gs *WorldScene) cleanupEntities() {
	activeBullets := make([]*entities.Bullet, 0, len(gs.bullets))
	for _, bullet := range gs.bullets {
		if !bullet.IsExpired() {
			activeBullets = append(activeBullets, bullet)
		} else {
			globals.Collision.UnregisterCollidable(bullet)
		}
	}
	gs.bullets = activeBullets

	activeEnemies := make([]*entities.Enemy, 0, len(gs.enemies))
	for _, enemy := range gs.enemies {
		if enemy.IsAlive() {
			activeEnemies = append(activeEnemies, enemy)
		} else {
			globals.Collision.UnregisterCollidable(enemy)
		}
	}
	gs.enemies = activeEnemies

	activePickups := make([]*entities.HealthPickup, 0, len(gs.healthPickups))
	for _, pickup := range gs.healthPickups {
		if pickup.IsActive() {
			activePickups = append(activePickups, pickup)
		} else {
			globals.Triggers.UnregisterTrigger(pickup)
		}
	}
	gs.healthPickups = activePickups
}

func floatToString(f float32, decimals int) string {
	if decimals == 0 {
		return fmt.Sprintf("%.0f", f)
	}
	return fmt.Sprintf("%.2f", f)
}

func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}
