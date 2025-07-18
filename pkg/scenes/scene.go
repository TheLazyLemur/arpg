package scenes

import (
	"arpg/pkg/config"
	"arpg/pkg/rendering"
)

type Scene interface {
	Initialize() error
	Update(deltaTime float32) error
	Render(renderer *rendering.Renderer) error
	HandleInput(deltaTime float32) error
	Cleanup() error
	GetName() string
	ShouldTransition() bool
	GetNextScene() string
}

type SceneManager struct {
	scenes        map[string]Scene
	currentScene  Scene
	nextSceneName string
	transitioning bool
	config        *config.Config
}

func NewSceneManager(cfg *config.Config) *SceneManager {
	return &SceneManager{
		scenes:        make(map[string]Scene),
		currentScene:  nil,
		nextSceneName: "",
		transitioning: false,
		config:        cfg,
	}
}

func (sm *SceneManager) RegisterScene(name string, scene Scene) {
	sm.scenes[name] = scene
}

func (sm *SceneManager) SetCurrentScene(name string) error {
	scene, exists := sm.scenes[name]
	if !exists {
		return &SceneError{Type: "scene_not_found", Message: "Scene not found: " + name}
	}

	if sm.currentScene != nil {
		if err := sm.currentScene.Cleanup(); err != nil {
			return &SceneError{
				Type:    "cleanup_failed",
				Message: "Failed to cleanup current scene: " + err.Error(),
			}
		}
	}

	if err := scene.Initialize(); err != nil {
		return &SceneError{
			Type:    "init_failed",
			Message: "Failed to initialize scene: " + err.Error(),
		}
	}

	sm.currentScene = scene
	sm.transitioning = false
	return nil
}

func (sm *SceneManager) Update(deltaTime float32) error {
	if sm.currentScene == nil {
		return &SceneError{Type: "no_scene", Message: "No current scene set"}
	}

	if err := sm.currentScene.Update(deltaTime); err != nil {
		return err
	}

	if err := sm.currentScene.HandleInput(deltaTime); err != nil {
		return err
	}

	if sm.currentScene.ShouldTransition() && !sm.transitioning {
		sm.nextSceneName = sm.currentScene.GetNextScene()
		sm.transitioning = true

		if err := sm.SetCurrentScene(sm.nextSceneName); err != nil {
			return err
		}
	}

	return nil
}

func (sm *SceneManager) Render(renderer *rendering.Renderer) error {
	if sm.currentScene == nil {
		return &SceneError{Type: "no_scene", Message: "No current scene set"}
	}

	return sm.currentScene.Render(renderer)
}

func (sm *SceneManager) GetCurrentScene() Scene {
	return sm.currentScene
}

func (sm *SceneManager) GetCurrentSceneName() string {
	if sm.currentScene == nil {
		return ""
	}
	return sm.currentScene.GetName()
}

func (sm *SceneManager) IsTransitioning() bool {
	return sm.transitioning
}

func (sm *SceneManager) Cleanup() error {
	if sm.currentScene != nil {
		return sm.currentScene.Cleanup()
	}
	return nil
}

type SceneError struct {
	Type    string
	Message string
}

func (e *SceneError) Error() string {
	return e.Message
}

