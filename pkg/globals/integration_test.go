package globals

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Test that all three systems (collision, trigger, input) can be initialized together
func TestSystemsInitialization(t *testing.T) {
	// given
	// ... no pre-existing systems

	// when
	// ... all three systems are initialized
	InitCollision()
	InitTriggers()
	InitInput()

	// then
	// ... should have a non-nil collision system
	// ... and should have a non-nil trigger system
	// ... and should have a non-nil input system
	if Collision == nil {
		t.Fatal("Collision system not initialized")
	}
	if Triggers == nil {
		t.Fatal("Trigger system not initialized")
	}
	if InputSystem == nil {
		t.Fatal("Input system not initialized")
	}
}

// Test that collision and trigger systems work together correctly
func TestCollisionTriggerIntegration(t *testing.T) {
	// given
	// ... initialized collision and trigger systems
	// ... a solid obstacle that blocks movement at position (5,0,0)
	// ... a trigger that doesn't block movement at position (0,0,0)
	// ... a player at position (0,0,0)
	// ... all entities registered with appropriate systems
	InitCollision()
	InitTriggers()
	obstacle := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 6, Y: 1, Z: 1},
		},
		CollisionTags: []string{"obstacle"},
		ActiveState:   true,
	}
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
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(obstacle)
	Collision.RegisterCollidable(player)
	Triggers.RegisterCollidable(player) // Player needs to be in trigger system
	Triggers.RegisterTrigger(healthPickup)

	// when
	// ... checking if player can move to trigger position
	triggerPos := rl.Vector3{X: 0.5, Y: 0.5, Z: 0.5}
	canMoveToTrigger := Collision.CheckMovement(player, triggerPos)

	// then
	// ... should allow movement (triggers don't block)
	if !canMoveToTrigger {
		t.Fatal("Player should be able to move to trigger position")
	}

	// when
	// ... the trigger system updates
	Triggers.Update()

	// then
	// ... should activate the trigger when player overlaps
	if !healthPickup.CallbackCalled {
		t.Fatal("Trigger should activate when player overlaps")
	}

	// when
	// ... checking if player can move to obstacle position
	obstaclePos := rl.Vector3{X: 5.5, Y: 0.5, Z: 0.5}
	canMoveToObstacle := Collision.CheckMovement(player, obstaclePos)

	// then
	// ... should block movement (obstacles block)
	if canMoveToObstacle {
		t.Fatal("Player should not be able to move to obstacle position")
	}
}

// Test that systems can be cleared without affecting each other
func TestSystemIndependence(t *testing.T) {
	// given
	// ... initialized collision and trigger systems
	// ... a collidable object
	// ... a trigger object
	// ... both objects registered with their respective systems
	InitCollision()
	InitTriggers()
	collidable := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	trigger := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	Collision.RegisterCollidable(collidable)
	Triggers.RegisterTrigger(trigger)

	// when
	// ... verifying initial state
	// then
	// ... collision system should have 1 object
	// ... and trigger system should have 1 object
	if len(Collision.collidables) != 1 {
		t.Fatal("Collision system should have 1 object")
	}
	if len(Triggers.triggers) != 1 {
		t.Fatal("Trigger system should have 1 object")
	}

	// when
	// ... the collision system is cleared
	Collision.ClearAll()

	// then
	// ... collision system should be empty
	// ... and trigger system should still have 1 object
	if len(Collision.collidables) != 0 {
		t.Fatal("Collision system should be cleared")
	}
	if len(Triggers.triggers) != 1 {
		t.Fatal("Trigger system should still have 1 object")
	}

	// when
	// ... the trigger system is cleared
	Triggers.ClearAll()

	// then
	// ... trigger system should be empty
	if len(Triggers.triggers) != 0 {
		t.Fatal("Trigger system should be cleared")
	}
}

// Test realistic game scenario: player moves through world with obstacles and pickups
func TestRealisticGameScenario(t *testing.T) {
	// given
	// ... initialized collision, trigger, and input systems
	// ... a player at position (0,0,0)
	// ... a solid wall obstacle at position (3,0,0)
	// ... an enemy at position (6,0,0)
	// ... a health pickup trigger at position (1,0,0)
	// ... all entities registered with appropriate systems
	InitCollision()
	InitTriggers()
	InitInput()
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	wall := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 3, Y: 0, Z: 0},
			Max: rl.Vector3{X: 4, Y: 1, Z: 1},
		},
		CollisionTags: []string{"obstacle"},
		ActiveState:   true,
	}
	enemy := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 6, Y: 0, Z: 0},
			Max: rl.Vector3{X: 7, Y: 1, Z: 1},
		},
		CollisionTags: []string{"enemy"},
		ActiveState:   true,
	}
	healthPickup := &MockTriggerable{
		TriggerBounds: rl.BoundingBox{
			Min: rl.Vector3{X: 1, Y: 0, Z: 0},
			Max: rl.Vector3{X: 2, Y: 1, Z: 1},
		},
		TriggerTags: []string{"health_pickup"},
		ActiveState: true,
	}
	Collision.RegisterCollidable(player)
	Triggers.RegisterCollidable(player) // Player needs to be in trigger system
	Collision.RegisterCollidable(wall)
	Collision.RegisterCollidable(enemy)
	Triggers.RegisterTrigger(healthPickup)

	// when
	// ... checking if player can move to empty space
	emptySpace := rl.Vector3{X: 1, Y: 0.5, Z: 0.5}
	canMoveToEmpty := Collision.CheckMovement(player, emptySpace)

	// then
	// ... should allow movement to empty space
	if !canMoveToEmpty {
		t.Fatal("Player should be able to move to empty space")
	}

	// when
	// ... checking if player can move through health pickup
	healthPickupPos := rl.Vector3{X: 1.5, Y: 0.5, Z: 0.5}
	canMoveToPickup := Collision.CheckMovement(player, healthPickupPos)

	// then
	// ... should allow movement (triggers don't block)
	if !canMoveToPickup {
		t.Fatal("Player should be able to move through health pickup")
	}

	// when
	// ... checking if player can move through wall
	wallPos := rl.Vector3{X: 3.5, Y: 0.5, Z: 0.5}
	canMoveToWall := Collision.CheckMovement(player, wallPos)

	// then
	// ... should block movement
	if canMoveToWall {
		t.Fatal("Player should not be able to move through wall")
	}

	// when
	// ... checking if player can move through enemy
	enemyPos := rl.Vector3{X: 6.5, Y: 0.5, Z: 0.5}
	canMoveToEnemy := Collision.CheckMovement(player, enemyPos)

	// then
	// ... should block movement
	if canMoveToEnemy {
		t.Fatal("Player should not be able to move through enemy")
	}

	// given
	// ... the player has moved to the health pickup position
	player.BoundingBox = rl.BoundingBox{
		Min: rl.Vector3{X: 1, Y: 0, Z: 0},
		Max: rl.Vector3{X: 2, Y: 1, Z: 1},
	}

	// when
	// ... the trigger system updates
	Triggers.Update()

	// then
	// ... should activate the health pickup trigger
	if !healthPickup.CallbackCalled {
		t.Fatal("Health pickup should trigger when player moves to it")
	}

	// given
	// ... the player has moved to the enemy position
	player.BoundingBox = rl.BoundingBox{
		Min: rl.Vector3{X: 6, Y: 0, Z: 0},
		Max: rl.Vector3{X: 7, Y: 1, Z: 1},
	}

	// when
	// ... the collision system updates
	Collision.Update()

	// then
	// ... should trigger collision callbacks for both player and enemy
	if !player.CallbackCalled {
		t.Fatal("Player should receive collision callback when overlapping with enemy")
	}
	if !enemy.CallbackCalled {
		t.Fatal("Enemy should receive collision callback when overlapping with player")
	}
}

// Test edge case: object that is both collidable and triggerable
func TestDualSystemObject(t *testing.T) {
	// given
	// ... initialized collision and trigger systems
	// ... an object that implements both Collidable and Triggerable interfaces
	// ... a player overlapping the dual object
	// ... both objects registered with both systems
	InitCollision()
	InitTriggers()
	dualObject := &DualSystemObject{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"obstacle"},
		TriggerTags:   []string{"health_pickup"},
		ActiveState:   true,
	}
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(dualObject)
	Triggers.RegisterTrigger(dualObject)
	Collision.RegisterCollidable(player)
	Triggers.RegisterCollidable(player)

	// when
	// ... both systems update
	Collision.Update()
	Triggers.Update()

	// then
	// ... should receive collision callback
	// ... and should receive trigger callback
	if !dualObject.CollisionCallbackCalled {
		t.Fatal("Dual object should receive collision callback")
	}
	if !dualObject.TriggerCallbackCalled {
		t.Fatal("Dual object should receive trigger callback")
	}
}

// DualSystemObject implements both Collidable and Triggerable for testing
type DualSystemObject struct {
	BoundingBox             rl.BoundingBox
	CollisionTags           []string
	TriggerTags             []string
	ActiveState             bool
	CollisionCallbackCalled bool
	TriggerCallbackCalled   bool
	CollisionCallbackOther  Collidable
	TriggerCallbackOther    Collidable
}

func (d *DualSystemObject) GetBoundingBox() rl.BoundingBox {
	return d.BoundingBox
}

func (d *DualSystemObject) GetCollisionTags() []string {
	return d.CollisionTags
}

func (d *DualSystemObject) OnCollision(other Collidable) {
	d.CollisionCallbackCalled = true
	d.CollisionCallbackOther = other
}

func (d *DualSystemObject) IsActive() bool {
	return d.ActiveState
}

func (d *DualSystemObject) GetTriggerBounds() rl.BoundingBox {
	return d.BoundingBox
}

func (d *DualSystemObject) GetTriggerTags() []string {
	return d.TriggerTags
}

func (d *DualSystemObject) OnTriggerEnter(other Collidable) {
	d.TriggerCallbackCalled = true
	d.TriggerCallbackOther = other
}

