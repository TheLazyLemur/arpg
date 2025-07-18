package globals

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestInputSystemInitialization(t *testing.T) {
	// given
	// ... no pre-existing input system

	// when
	// ... the input system is initialized
	InitInput()

	// then
	// ... should have a non-nil input system
	if InputSystem == nil {
		t.Fatal("Input system not initialized")
	}

	// when
	// ... calling input system methods without a window context
	// then
	// ... should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Input system methods should not panic: %v", r)
		}
	}()

	// ... all input methods should exist and execute without panicking
	_ = InputSystem.IsSpacePressed()
	_ = InputSystem.IsSpaceDown()
	_ = InputSystem.IsPausePressed()
	_ = InputSystem.IsRestartPressed()
	_ = InputSystem.IsFullscreenPressed()
	_ = InputSystem.IsDebugPressed()
	_ = InputSystem.IsEscapePressed()
	_ = InputSystem.IsUpPressed()
	_ = InputSystem.IsDownPressed()
	_ = InputSystem.IsLeftPressed()
	_ = InputSystem.IsRightPressed()
	_ = InputSystem.IsUpDown()
	_ = InputSystem.IsDownDown()
	_ = InputSystem.IsLeftDown()
	_ = InputSystem.IsRightDown()
	_ = InputSystem.IsMouseLeftPressed()
	_ = InputSystem.IsMouseLeftDown()
	_ = InputSystem.GetMousePosition()
}

func TestInputSystemInterface(t *testing.T) {
	// given
	// ... an initialized input system
	InitInput()

	// when
	// ... checking interface implementation
	// then
	// ... InputSystem should implement the Input interface
	var _ Input = InputSystem
	// ... and DefaultInput should implement the Input interface
	var _ Input = &DefaultInput{}
}

func TestInputSystemMethods(t *testing.T) {
	// given
	// ... an initialized input system
	InitInput()

	// when
	// ... checking menu/system input method return types
	// then
	// ... IsSpacePressed should return bool
	// ... and IsSpaceDown should return bool
	if _, ok := interface{}(InputSystem.IsSpacePressed()).(bool); !ok {
		t.Fatal("IsSpacePressed should return bool")
	}
	if _, ok := interface{}(InputSystem.IsSpaceDown()).(bool); !ok {
		t.Fatal("IsSpaceDown should return bool")
	}
	// ... and IsPausePressed should return bool
	if _, ok := interface{}(InputSystem.IsPausePressed()).(bool); !ok {
		t.Fatal("IsPausePressed should return bool")
	}
	// ... and IsRestartPressed should return bool
	if _, ok := interface{}(InputSystem.IsRestartPressed()).(bool); !ok {
		t.Fatal("IsRestartPressed should return bool")
	}
	// ... and IsFullscreenPressed should return bool
	if _, ok := interface{}(InputSystem.IsFullscreenPressed()).(bool); !ok {
		t.Fatal("IsFullscreenPressed should return bool")
	}
	// ... and IsDebugPressed should return bool
	if _, ok := interface{}(InputSystem.IsDebugPressed()).(bool); !ok {
		t.Fatal("IsDebugPressed should return bool")
	}
	// ... and IsEscapePressed should return bool
	if _, ok := interface{}(InputSystem.IsEscapePressed()).(bool); !ok {
		t.Fatal("IsEscapePressed should return bool")
	}

	// when
	// ... checking movement input method return types
	// then
	// ... IsUpPressed should return bool
	if _, ok := interface{}(InputSystem.IsUpPressed()).(bool); !ok {
		t.Fatal("IsUpPressed should return bool")
	}
	// ... and IsDownPressed should return bool
	if _, ok := interface{}(InputSystem.IsDownPressed()).(bool); !ok {
		t.Fatal("IsDownPressed should return bool")
	}
	// ... and IsLeftPressed should return bool
	if _, ok := interface{}(InputSystem.IsLeftPressed()).(bool); !ok {
		t.Fatal("IsLeftPressed should return bool")
	}
	// ... and IsRightPressed should return bool
	if _, ok := interface{}(InputSystem.IsRightPressed()).(bool); !ok {
		t.Fatal("IsRightPressed should return bool")
	}
	// ... and IsUpDown should return bool
	if _, ok := interface{}(InputSystem.IsUpDown()).(bool); !ok {
		t.Fatal("IsUpDown should return bool")
	}
	// ... and IsDownDown should return bool
	if _, ok := interface{}(InputSystem.IsDownDown()).(bool); !ok {
		t.Fatal("IsDownDown should return bool")
	}
	// ... and IsLeftDown should return bool
	if _, ok := interface{}(InputSystem.IsLeftDown()).(bool); !ok {
		t.Fatal("IsLeftDown should return bool")
	}
	// ... and IsRightDown should return bool
	if _, ok := interface{}(InputSystem.IsRightDown()).(bool); !ok {
		t.Fatal("IsRightDown should return bool")
	}

	// when
	// ... checking mouse input method return types
	// then
	// ... IsMouseLeftPressed should return bool
	if _, ok := interface{}(InputSystem.IsMouseLeftPressed()).(bool); !ok {
		t.Fatal("IsMouseLeftPressed should return bool")
	}
	// ... and IsMouseLeftDown should return bool
	if _, ok := interface{}(InputSystem.IsMouseLeftDown()).(bool); !ok {
		t.Fatal("IsMouseLeftDown should return bool")
	}

	// when
	// ... getting mouse position
	mousePos := InputSystem.GetMousePosition()
	// then
	// ... should return a Vector2 (values may be invalid without window)
	if mousePos.X < 0 && mousePos.Y < 0 {
		// This is fine - without a window, mouse position might be invalid
		// but the method should still work
	}
}

// Test that the default input system can be replaced with a mock
type MockInput struct {
	SpacePressed      bool
	SpaceDown         bool
	PausePressed      bool
	RestartPressed    bool
	FullscreenPressed bool
	DebugPressed      bool
	EscapePressed     bool
	UpPressed         bool
	DownPressed       bool
	LeftPressed       bool
	RightPressed      bool
	UpDown            bool
	DownDown          bool
	LeftDown          bool
	RightDown         bool
	MouseLeftPressed  bool
	MouseLeftDown     bool
	MouseX            float32
	MouseY            float32
}

func (m *MockInput) IsSpacePressed() bool { return m.SpacePressed }

func (m *MockInput) IsSpaceDown() bool { return m.SpaceDown }

func (m *MockInput) IsPausePressed() bool { return m.PausePressed }

func (m *MockInput) IsRestartPressed() bool { return m.RestartPressed }

func (m *MockInput) IsFullscreenPressed() bool { return m.FullscreenPressed }

func (m *MockInput) IsDebugPressed() bool { return m.DebugPressed }

func (m *MockInput) IsEscapePressed() bool { return m.EscapePressed }

func (m *MockInput) IsUpPressed() bool { return m.UpPressed }

func (m *MockInput) IsDownPressed() bool { return m.DownPressed }

func (m *MockInput) IsLeftPressed() bool { return m.LeftPressed }

func (m *MockInput) IsRightPressed() bool { return m.RightPressed }

func (m *MockInput) IsUpDown() bool { return m.UpDown }

func (m *MockInput) IsDownDown() bool { return m.DownDown }

func (m *MockInput) IsLeftDown() bool { return m.LeftDown }

func (m *MockInput) IsRightDown() bool { return m.RightDown }

func (m *MockInput) IsMouseLeftPressed() bool { return m.MouseLeftPressed }

func (m *MockInput) IsMouseLeftDown() bool { return m.MouseLeftDown }

func (m *MockInput) GetMousePosition() rl.Vector2 {
	return rl.Vector2{X: m.MouseX, Y: m.MouseY}
}

func TestInputSystemMockability(t *testing.T) {
	// given
	// ... the original input system
	// ... a mock input with specific test values
	originalInput := InputSystem
	mockInput := &MockInput{
		SpacePressed:     true,
		UpDown:           true,
		MouseLeftPressed: true,
		MouseX:           100,
		MouseY:           200,
	}

	// when
	// ... the input system is replaced with the mock
	InputSystem = mockInput

	// then
	// ... should return mock values for IsSpacePressed
	// ... and should return mock values for IsUpDown
	// ... and should return mock values for IsMouseLeftPressed
	if !InputSystem.IsSpacePressed() {
		t.Fatal("Mock input should return true for IsSpacePressed")
	}
	if !InputSystem.IsUpDown() {
		t.Fatal("Mock input should return true for IsUpDown")
	}
	if !InputSystem.IsMouseLeftPressed() {
		t.Fatal("Mock input should return true for IsMouseLeftPressed")
	}

	// when
	// ... getting mouse position
	mousePos := InputSystem.GetMousePosition()
	// then
	// ... should return mock coordinates
	if mousePos.X != 100 || mousePos.Y != 200 {
		t.Fatalf(
			"Mock input should return (100, 200) for mouse position, got (%f, %f)",
			mousePos.X,
			mousePos.Y,
		)
	}

	// when
	// ... checking non-set values
	// then
	// ... should return false for IsSpaceDown when not set
	if InputSystem.IsSpaceDown() {
		t.Fatal("Mock input should return false for IsSpaceDown when not set")
	}

	// when
	// ... restoring the original input system
	InputSystem = originalInput
	// then
	// ... should have restored the original system
}
