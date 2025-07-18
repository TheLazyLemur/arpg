package globals

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Input interface {
	// Menu/System inputs
	IsSpacePressed() bool
	IsSpaceDown() bool
	IsPausePressed() bool
	IsRestartPressed() bool
	IsFullscreenPressed() bool
	IsDebugPressed() bool
	IsEscapePressed() bool

	// Movement inputs
	IsUpPressed() bool
	IsDownPressed() bool
	IsLeftPressed() bool
	IsRightPressed() bool
	IsUpDown() bool
	IsDownDown() bool
	IsLeftDown() bool
	IsRightDown() bool

	// Mouse inputs
	IsMouseLeftPressed() bool
	IsMouseLeftDown() bool
	GetMousePosition() rl.Vector2
}

var InputSystem Input

type DefaultInput struct{}

func InitInput() {
	InputSystem = &DefaultInput{}
}

// Menu/System inputs
func (di *DefaultInput) IsSpacePressed() bool {
	return rl.IsKeyPressed(rl.KeySpace)
}

func (di *DefaultInput) IsSpaceDown() bool {
	return rl.IsKeyDown(rl.KeySpace)
}

func (di *DefaultInput) IsPausePressed() bool {
	return rl.IsKeyPressed(rl.KeyP)
}

func (di *DefaultInput) IsRestartPressed() bool {
	return rl.IsKeyPressed(rl.KeyR)
}

func (di *DefaultInput) IsFullscreenPressed() bool {
	return rl.IsKeyPressed(rl.KeyF11)
}

func (di *DefaultInput) IsDebugPressed() bool {
	return rl.IsKeyPressed(rl.KeyF3)
}

func (di *DefaultInput) IsEscapePressed() bool {
	return rl.IsKeyPressed(rl.KeyEscape)
}

// Movement inputs - pressed
func (di *DefaultInput) IsUpPressed() bool {
	return rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyUp)
}

func (di *DefaultInput) IsDownPressed() bool {
	return rl.IsKeyPressed(rl.KeyS) || rl.IsKeyPressed(rl.KeyDown)
}

func (di *DefaultInput) IsLeftPressed() bool {
	return rl.IsKeyPressed(rl.KeyA) || rl.IsKeyPressed(rl.KeyLeft)
}

func (di *DefaultInput) IsRightPressed() bool {
	return rl.IsKeyPressed(rl.KeyD) || rl.IsKeyPressed(rl.KeyRight)
}

// Movement inputs - held down
func (di *DefaultInput) IsUpDown() bool {
	return rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp)
}

func (di *DefaultInput) IsDownDown() bool {
	return rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown)
}

func (di *DefaultInput) IsLeftDown() bool {
	return rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)
}

func (di *DefaultInput) IsRightDown() bool {
	return rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)
}

// Mouse inputs
func (di *DefaultInput) IsMouseLeftPressed() bool {
	return rl.IsMouseButtonPressed(rl.MouseLeftButton)
}

func (di *DefaultInput) IsMouseLeftDown() bool {
	return rl.IsMouseButtonDown(rl.MouseLeftButton)
}

func (di *DefaultInput) GetMousePosition() rl.Vector2 {
	return rl.GetMousePosition()
}
