package game

import (
	"fmt"
	"log"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	gameScenes "arpg/internal/scenes"
	"arpg/pkg/config"
	"arpg/pkg/rendering"
	"arpg/pkg/scenes"
)

type Game struct {
	config       *config.Config
	renderer     *rendering.Renderer
	sceneManager *scenes.SceneManager
	running      bool
}

func New(cfg *config.Config) *Game {
	g := &Game{
		config:  cfg,
		running: true,
	}

	g.renderer = rendering.NewRenderer(cfg)
	g.sceneManager = scenes.NewSceneManager(cfg)

	return g
}

func (g *Game) Initialize() error {
	g.renderer.Initialize()

	menuScene := gameScenes.NewMenuScene(g.config)
	gameScene := gameScenes.NewGameScene(g.config)

	g.sceneManager.RegisterScene("menu", menuScene)
	g.sceneManager.RegisterScene("game", gameScene)

	if err := g.sceneManager.SetCurrentScene("game"); err != nil {
		return fmt.Errorf("failed to set initial scene: %w", err)
	}

	log.Println("Game initialized successfully")
	return nil
}

func (g *Game) Run() error {
	if err := g.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize game: %w", err)
	}

	defer g.Shutdown()

	for g.running && !g.renderer.ShouldClose() {
		if err := g.Update(); err != nil {
			if isExitError(err) {
				log.Println("Exit requested by user")
				g.running = false
				continue
			}
			return fmt.Errorf("game update error: %w", err)
		}

		if err := g.Render(); err != nil {
			return fmt.Errorf("game render error: %w", err)
		}
	}

	return nil
}

func (g *Game) Update() error {
	deltaTime := rl.GetFrameTime()

	if err := g.sceneManager.Update(deltaTime); err != nil {
		return err
	}

	return nil
}

func (g *Game) Render() error {
	if err := g.sceneManager.Render(g.renderer); err != nil {
		return err
	}

	return nil
}

func (g *Game) Shutdown() {
	if err := g.sceneManager.Cleanup(); err != nil {
		log.Printf("Error during scene cleanup: %v", err)
	}

	g.renderer.Shutdown()
	log.Println("Game shutdown complete")
}

func (g *Game) GetConfig() *config.Config {
	return g.config
}

func (g *Game) IsRunning() bool {
	return g.running
}

func (g *Game) Stop() {
	g.running = false
}

func (g *Game) GetCurrentScene() string {
	return g.sceneManager.GetCurrentSceneName()
}

func isExitError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "exit_requested") ||
		strings.Contains(errMsg, "user requested exit")
}
