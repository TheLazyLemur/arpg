package globals

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Mock triggerable for testing
type MockTriggerable struct {
	TriggerBounds  rl.BoundingBox
	TriggerTags    []string
	ActiveState    bool
	CallbackCalled bool
	CallbackOther  Collidable
	CallbackCount  int
}

func (m *MockTriggerable) GetTriggerBounds() rl.BoundingBox {
	return m.TriggerBounds
}

func (m *MockTriggerable) GetTriggerTags() []string {
	return m.TriggerTags
}

func (m *MockTriggerable) OnTriggerEnter(other Collidable) {
	m.CallbackCalled = true
	m.CallbackOther = other
	m.CallbackCount++
}

func (m *MockTriggerable) IsActive() bool {
	return m.ActiveState
}

func TestTriggerSystemInitialization(t *testing.T) {
	// given
	// ... no pre-existing trigger system
	// when
	// ... the trigger system is initialized
	InitTriggers()
	// then
	// ... should have a non-nil trigger system
	// ... and should have an initialized triggers slice
	// ... and should have zero triggers initially
	if Triggers == nil {
		t.Fatal("Trigger system not initialized")
	}
	if Triggers.triggers == nil {
		t.Fatal("Triggers slice not initialized")
	}
	if len(Triggers.triggers) != 0 {
		t.Fatal("Expected empty triggers slice on initialization")
	}
}

func TestTriggerSystemRegistration(t *testing.T) {
	// given
	// ... an initialized trigger system
	// ... a mock trigger with health pickup properties
	InitTriggers()
	mockTrigger := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	// when
	// ... the trigger is registered
	Triggers.RegisterTrigger(mockTrigger)
	// then
	// ... should have exactly one trigger
	// ... and should be the registered mock trigger
	if len(Triggers.triggers) != 1 {
		t.Fatalf("Expected 1 trigger, got %d", len(Triggers.triggers))
	}
	if Triggers.triggers[0] != mockTrigger {
		t.Fatal("Registered trigger not found in system")
	}
}

func TestTriggerSystemUnregistration(t *testing.T) {
	// given
	// ... an initialized trigger system
	// ... two mock triggers registered in the system
	InitTriggers()
	trigger1 := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	trigger2 := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 2, Y: 0, Z: 0},
			Max: rl.Vector3{X: 3, Y: 1, Z: 1},
		},
		TriggerTags: []string{"damage_zone"},
		ActiveState: true,
	}
	Triggers.RegisterTrigger(trigger1)
	Triggers.RegisterTrigger(trigger2)
	// when
	// ... the first trigger is unregistered
	Triggers.UnregisterTrigger(trigger1)
	// then
	// ... should have exactly one trigger remaining
	// ... and should be the second trigger
	if len(Triggers.triggers) != 1 {
		t.Fatalf("Expected 1 trigger after unregistration, got %d", len(Triggers.triggers))
	}
	if Triggers.triggers[0] != trigger2 {
		t.Fatal("Wrong trigger remained after unregistration")
	}
}

func TestTriggerSystemClearAll(t *testing.T) {
	// given
	// ... an initialized trigger system
	// ... two triggers registered in the system
	InitTriggers()
	trigger1 := &MockTriggerable{TriggerTags: []string{"health_pickup"}, ActiveState: true}
	trigger2 := &MockTriggerable{TriggerTags: []string{"damage_zone"}, ActiveState: true}
	Triggers.RegisterTrigger(trigger1)
	Triggers.RegisterTrigger(trigger2)
	// when
	// ... the clear all method is called
	Triggers.ClearAll()
	// then
	// ... should have no triggers remaining
	if len(Triggers.triggers) != 0 {
		t.Fatalf("Expected 0 triggers after ClearAll, got %d", len(Triggers.triggers))
	}
}

func TestTriggerRules(t *testing.T) {
	// given
	// ... an initialized trigger system
	// ... a health pickup trigger
	// ... a player collidable
	InitTriggers()
	healthPickup := &MockTriggerable{
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	player := &MockCollidable{
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	// when
	// ... checking if health pickup should trigger on player
	playerTriggers := Triggers.shouldTrigger(healthPickup, player)

	// then
	// ... should allow trigger
	if !playerTriggers {
		t.Fatal("Health pickup should trigger on player")
	}
	// given
	// ... an enemy collidable
	enemy := &MockCollidable{
		CollisionTags: []string{"enemy"},
		ActiveState:   true,
	}
	// when
	// ... checking if health pickup should trigger on enemy
	enemyTriggers := Triggers.shouldTrigger(healthPickup, enemy)

	// then
	// ... should not allow trigger
	if enemyTriggers {
		t.Fatal("Health pickup should NOT trigger on enemy")
	}
	// given
	// ... a trigger with unknown tag
	unknownTrigger := &MockTriggerable{
		TriggerTags: []string{"unknown_trigger"},
		ActiveState: true,
	}
	// when
	// ... checking if unknown trigger should trigger on player
	unknownTriggers := Triggers.shouldTrigger(unknownTrigger, player)
	// then
	// ... should not allow trigger
	if unknownTriggers {
		t.Fatal("Unknown trigger should not trigger on player")
	}
}

func TestTriggerActivationWithCollisionSystem(t *testing.T) {
	// given
	// ... initialized trigger and collision systems
	// ... a health pickup trigger at position (0,0,0)
	// ... a player overlapping the trigger
	// ... trigger and player registered with respective systems
	InitTriggers()
	InitCollision()
	healthPickup := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	Triggers.RegisterTrigger(healthPickup)
	Collision.RegisterCollidable(player)
	Triggers.RegisterCollidable(player) // Player needs to be registered with trigger system too
	// when
	// ... the trigger system updates
	Triggers.Update()
	// then
	// ... should activate the health pickup trigger
	if !healthPickup.CallbackCalled {
		t.Fatal("Health pickup trigger should have been activated by player")
	}
	// ... and should pass the player as the callback other
	if healthPickup.CallbackOther != player {
		t.Fatal("Health pickup should have received player as callback other")
	}
}

func TestTriggerNoActivationWhenNotOverlapping(t *testing.T) {
	// given
	// ... initialized trigger and collision systems
	// ... a trigger at position (0,0,0)
	// ... a player far away at position (5,0,0)
	// ... both registered with their respective systems
	InitTriggers()
	InitCollision()
	trigger := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 6, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	Triggers.RegisterTrigger(trigger)
	Collision.RegisterCollidable(player)
	Triggers.RegisterCollidable(player)

	// when
	// ... the trigger system updates
	Triggers.Update()

	// then
	// ... should not activate the trigger
	if trigger.CallbackCalled {
		t.Fatal("Trigger should not activate when not overlapping")
	}
}

func TestTriggerInactiveObjects(t *testing.T) {
	// given
	// ... initialized trigger and collision systems
	// ... an inactive trigger at position (0,0,0)
	// ... an active player overlapping the trigger
	// ... both registered with their respective systems
	InitTriggers()
	InitCollision()
	inactiveTrigger := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: false, // Inactive
	}
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	Triggers.RegisterTrigger(inactiveTrigger)
	Collision.RegisterCollidable(player)
	Triggers.RegisterCollidable(player)

	// when
	// ... the trigger system updates
	Triggers.Update()

	// then
	// ... should not activate the inactive trigger
	if inactiveTrigger.CallbackCalled {
		t.Fatal("Inactive trigger should not activate")
	}

	// given
	// ... an active trigger at position (0,0,0)
	// ... an inactive player overlapping the trigger
	// ... a fresh system with the active trigger and inactive player
	activeTrigger := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	inactivePlayer := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   false, // Inactive
	}
	Triggers.ClearAll()
	Collision.ClearAll()
	Triggers.RegisterTrigger(activeTrigger)
	Collision.RegisterCollidable(inactivePlayer)
	Triggers.RegisterCollidable(inactivePlayer)

	// when
	// ... the trigger system updates
	Triggers.Update()

	// then
	// ... should not activate the trigger with inactive collidable
	if activeTrigger.CallbackCalled {
		t.Fatal("Trigger should not activate with inactive collidable")
	}
}

func TestTriggerMultipleUpdates(t *testing.T) {
	// given
	// ... initialized trigger and collision systems
	// ... an overlapping trigger and collidable
	// ... a player overlapping the trigger
	// ... both registered with their respective systems
	InitTriggers()
	InitCollision()
	trigger := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	Triggers.RegisterTrigger(trigger)
	Collision.RegisterCollidable(player)
	Triggers.RegisterCollidable(player)

	// when
	// ... the trigger system updates multiple times
	Triggers.Update()
	Triggers.Update()
	Triggers.Update()

	// then
	// ... should trigger multiple times (continuous trigger)
	if trigger.CallbackCount != 3 {
		t.Fatalf("Expected trigger to be called 3 times, got %d", trigger.CallbackCount)
	}
}

func TestTriggerCollisionDetection(t *testing.T) {
	// given
	// ... an initialized trigger system
	// ... two identical bounding boxes
	InitTriggers()
	bounds1 := rl.BoundingBox{
		Min: rl.Vector3{X: 0, Y: 0, Z: 0},
		Max: rl.Vector3{X: 1, Y: 1, Z: 1},
	}
	bounds2 := rl.BoundingBox{
		Min: rl.Vector3{X: 0, Y: 0, Z: 0},
		Max: rl.Vector3{X: 1, Y: 1, Z: 1},
	}

	// when
	// ... checking collision between identical boxes
	exactOverlap := Triggers.checkTriggerCollision(bounds1, bounds2)

	// then
	// ... should detect collision
	if !exactOverlap {
		t.Fatal("Exact overlap should trigger collision")
	}

	// given
	// ... a partially overlapping bounding box
	bounds3 := rl.BoundingBox{
		Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
		Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
	}

	// when
	// ... checking collision with partial overlap
	partialOverlap := Triggers.checkTriggerCollision(bounds1, bounds3)

	// then
	// ... should detect collision
	if !partialOverlap {
		t.Fatal("Partial overlap should trigger collision")
	}

	// given
	// ... a non-overlapping bounding box
	bounds4 := rl.BoundingBox{
		Min: rl.Vector3{X: 2, Y: 0, Z: 0},
		Max: rl.Vector3{X: 3, Y: 1, Z: 1},
	}

	// when
	// ... checking collision with no overlap
	noOverlap := Triggers.checkTriggerCollision(bounds1, bounds4)

	// then
	// ... should not detect collision
	if noOverlap {
		t.Fatal("No overlap should not trigger collision")
	}
}

func TestTriggerSystemEdgeCases(t *testing.T) {
	// given
	// ... initialized trigger and collision systems
	InitTriggers()
	InitCollision()

	// when
	// ... the trigger system updates with empty systems
	Triggers.Update()

	// then
	// ... should not crash

	// given
	// ... a trigger with no collidables in the system
	trigger := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	// ... registered with the trigger system
	Triggers.RegisterTrigger(trigger)

	// when
	// ... the trigger system updates
	Triggers.Update()

	// then
	// ... should not activate the trigger
	if trigger.CallbackCalled {
		t.Fatal("Trigger should not activate with no collidables")
	}

	// given
	// ... a second trigger at the same position
	trigger2 := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	// ... a player overlapping both triggers
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	// ... both triggers and player registered with their systems
	Triggers.RegisterTrigger(trigger2)
	Collision.RegisterCollidable(player)
	Triggers.RegisterCollidable(player)

	// when
	// ... the trigger system updates
	Triggers.Update()

	// then
	// ... should activate both triggers with the same collidable
	if !trigger.CallbackCalled || !trigger2.CallbackCalled {
		t.Fatal("Both triggers should activate with same collidable")
	}
}
