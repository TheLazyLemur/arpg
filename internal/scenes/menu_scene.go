package scenes

import (
	"arpg/pkg/rendering"
	"arpg/pkg/config"
	
	rl "github.com/gen2brain/raylib-go/raylib"
)

// MenuScene represents the main menu
type MenuScene struct {
	config            *config.Config
	shouldTransition  bool
	nextScene         string
	
	// UI elements
	startButtonRect   rl.Rectangle
	startButtonHover  bool
	title             string
	
	// Colors
	backgroundColor   rl.Color
	titleColor        rl.Color
	buttonColor       rl.Color
	buttonHoverColor  rl.Color
	buttonTextColor   rl.Color
}

// NewMenuScene creates a new menu scene
func NewMenuScene(cfg *config.Config) *MenuScene {
	return &MenuScene{
		config:           cfg,
		shouldTransition: false,
		nextScene:        "",
		title:            "ARPG - 3D Action RPG",
		backgroundColor:  rl.NewColor(30, 30, 50, 255),
		titleColor:       rl.White,
		buttonColor:      rl.NewColor(70, 70, 120, 255),
		buttonHoverColor: rl.NewColor(100, 100, 160, 255),
		buttonTextColor:  rl.White,
	}
}

// Initialize sets up the menu scene
func (ms *MenuScene) Initialize() error {
	// Calculate button position (centered)
	buttonWidth := float32(200)
	buttonHeight := float32(60)
	centerX := float32(ms.config.Window.Width) / 2
	centerY := float32(ms.config.Window.Height) / 2
	
	ms.startButtonRect = rl.Rectangle{
		X:      centerX - buttonWidth/2,
		Y:      centerY - buttonHeight/2,
		Width:  buttonWidth,
		Height: buttonHeight,
	}
	
	ms.shouldTransition = false
	ms.nextScene = ""
	
	return nil
}

// Update handles menu logic
func (ms *MenuScene) Update(deltaTime float32) error {
	// Check if mouse is hovering over start button
	mousePos := rl.GetMousePosition()
	ms.startButtonHover = rl.CheckCollisionPointRec(mousePos, ms.startButtonRect)
	
	return nil
}

// Render draws the menu scene
func (ms *MenuScene) Render(renderer *rendering.Renderer) error {
	renderer.BeginFrame()
	
	// Draw background
	rl.ClearBackground(ms.backgroundColor)
	
	// Draw title
	titleFontSize := int32(48)
	titleText := ms.title
	titleWidth := rl.MeasureText(titleText, titleFontSize)
	titleX := (ms.config.Window.Width - titleWidth) / 2
	titleY := ms.config.Window.Height / 4
	
	rl.DrawText(titleText, titleX, titleY, titleFontSize, ms.titleColor)
	
	// Draw start button
	buttonColor := ms.buttonColor
	if ms.startButtonHover {
		buttonColor = ms.buttonHoverColor
	}
	
	rl.DrawRectangleRec(ms.startButtonRect, buttonColor)
	rl.DrawRectangleLinesEx(ms.startButtonRect, 2, rl.White)
	
	// Draw button text
	buttonText := "START GAME"
	buttonFontSize := int32(24)
	buttonTextWidth := rl.MeasureText(buttonText, buttonFontSize)
	buttonTextX := int32(ms.startButtonRect.X + (ms.startButtonRect.Width-float32(buttonTextWidth))/2)
	buttonTextY := int32(ms.startButtonRect.Y + (ms.startButtonRect.Height-float32(buttonFontSize))/2)
	
	rl.DrawText(buttonText, buttonTextX, buttonTextY, buttonFontSize, ms.buttonTextColor)
	
	// Draw instructions
	instructionText := "Use WASD to move, mouse to aim, left click to shoot"
	instructionFontSize := int32(16)
	instructionWidth := rl.MeasureText(instructionText, instructionFontSize)
	instructionX := (ms.config.Window.Width - instructionWidth) / 2
	instructionY := ms.config.Window.Height - 100
	
	rl.DrawText(instructionText, instructionX, instructionY, instructionFontSize, rl.Gray)
	
	// Draw additional info
	infoText := "Press ESC to exit game"
	infoFontSize := int32(14)
	infoWidth := rl.MeasureText(infoText, infoFontSize)
	infoX := (ms.config.Window.Width - infoWidth) / 2
	infoY := ms.config.Window.Height - 60
	
	rl.DrawText(infoText, infoX, infoY, infoFontSize, rl.DarkGray)
	
	renderer.EndFrame()
	return nil
}

// HandleInput processes input for the menu scene
func (ms *MenuScene) HandleInput(deltaTime float32) error {
	// Handle mouse click on start button
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && ms.startButtonHover {
		ms.shouldTransition = true
		ms.nextScene = "game"
	}
	
	// Handle Enter key to start game
	if rl.IsKeyPressed(rl.KeyEnter) {
		ms.shouldTransition = true
		ms.nextScene = "game"
	}
	
	// Handle Escape key to exit
	if rl.IsKeyPressed(rl.KeyEscape) {
		// Exit the application
		return &MenuError{Type: "exit_requested", Message: "User requested exit"}
	}
	
	return nil
}

// Cleanup releases menu scene resources
func (ms *MenuScene) Cleanup() error {
	// No resources to cleanup for menu scene
	return nil
}

// GetName returns the scene name
func (ms *MenuScene) GetName() string {
	return "menu"
}

// ShouldTransition returns true if scene wants to transition
func (ms *MenuScene) ShouldTransition() bool {
	return ms.shouldTransition
}

// GetNextScene returns the next scene to transition to
func (ms *MenuScene) GetNextScene() string {
	return ms.nextScene
}

// MenuError represents menu-related errors
type MenuError struct {
	Type    string
	Message string
}

func (e *MenuError) Error() string {
	return e.Message
}