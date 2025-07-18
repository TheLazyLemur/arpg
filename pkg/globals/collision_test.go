package globals

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Mock collidable for testing
type MockCollidable struct {
	BoundingBox    rl.BoundingBox
	CollisionTags  []string
	ActiveState    bool
	CallbackCalled bool
	CallbackOther  Collidable
}

func (m *MockCollidable) GetBoundingBox() rl.BoundingBox {
	return m.BoundingBox
}

func (m *MockCollidable) GetCollisionTags() []string {
	return m.CollisionTags
}

func (m *MockCollidable) OnCollision(other Collidable) {
	m.CallbackCalled = true
	m.CallbackOther = other
}

func (m *MockCollidable) IsActive() bool {
	return m.ActiveState
}

func TestCollisionSystemInitialization(t *testing.T) {
	// given
	// ... no pre-existing collision system
	// when
	// ... the collision system is initialized
	InitCollision()
	// then
	// ... should have a non-nil collision system
	// ... and should have an initialized collidables slice
	// ... and should have zero collidables initially
	if Collision == nil {
		t.Fatal("Collision system not initialized")
	}
	if Collision.collidables == nil {
		t.Fatal("Collidables slice not initialized")
	}
	if len(Collision.collidables) != 0 {
		t.Fatal("Expected empty collidables slice on initialization")
	}
}

func TestCollisionSystemRegistration(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... a mock collidable with test properties
	InitCollision()
	mock := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"test"},
		ActiveState:   true,
	}
	// when
	// ... the collidable is registered
	Collision.RegisterCollidable(mock)
	// then
	// ... should have exactly one collidable
	// ... and should be the registered mock
	if len(Collision.collidables) != 1 {
		t.Fatalf("Expected 1 collidable, got %d", len(Collision.collidables))
	}
	if Collision.collidables[0] != mock {
		t.Fatal("Registered collidable not found in system")
	}
	// given
	// ... a second mock collidable
	mock2 := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 2, Y: 0, Z: 0},
			Max: rl.Vector3{X: 3, Y: 1, Z: 1},
		},
		CollisionTags: []string{"test2"},
		ActiveState:   true,
	}
	// when
	// ... the second collidable is registered
	Collision.RegisterCollidable(mock2)
	// then
	// ... should have exactly two collidables
	if len(Collision.collidables) != 2 {
		t.Fatalf("Expected 2 collidables, got %d", len(Collision.collidables))
	}
}

func TestCollisionSystemUnregistration(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... two mock collidables registered in the system
	InitCollision()
	mock1 := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"test1"},
		ActiveState:   true,
	}
	mock2 := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 2, Y: 0, Z: 0},
			Max: rl.Vector3{X: 3, Y: 1, Z: 1},
		},
		CollisionTags: []string{"test2"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(mock1)
	Collision.RegisterCollidable(mock2)
	// when
	// ... the first collidable is unregistered
	Collision.UnregisterCollidable(mock1)
	// then
	// ... should have exactly one collidable remaining
	// ... and should be the second mock
	if len(Collision.collidables) != 1 {
		t.Fatalf("Expected 1 collidable after unregistration, got %d", len(Collision.collidables))
	}
	if Collision.collidables[0] != mock2 {
		t.Fatal("Wrong collidable remained after unregistration")
	}
	// given
	// ... a non-existent collidable
	mock3 := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 4, Y: 0, Z: 0},
			Max: rl.Vector3{X: 5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"test3"},
		ActiveState:   true,
	}
	// when
	// ... the non-existent collidable is unregistered
	Collision.UnregisterCollidable(mock3)
	// then
	// ... should still have exactly one collidable
	if len(Collision.collidables) != 1 {
		t.Fatalf(
			"Expected 1 collidable after unregistering non-existent, got %d",
			len(Collision.collidables),
		)
	}
}

func TestCollisionSystemClearAll(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... two collidables registered in the system
	InitCollision()
	mock1 := &MockCollidable{CollisionTags: []string{"test1"}, ActiveState: true}
	mock2 := &MockCollidable{CollisionTags: []string{"test2"}, ActiveState: true}
	Collision.RegisterCollidable(mock1)
	Collision.RegisterCollidable(mock2)
	// when
	// ... the clear all method is called
	Collision.ClearAll()
	// then
	// ... should have no collidables remaining
	if len(Collision.collidables) != 0 {
		t.Fatalf("Expected 0 collidables after ClearAll, got %d", len(Collision.collidables))
	}
}

func TestCollisionRules(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... a player collidable
	// ... an enemy collidable
	InitCollision()
	player := &MockCollidable{
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	enemy := &MockCollidable{
		CollisionTags: []string{"enemy"},
		ActiveState:   true,
	}
	// when
	// ... checking if player should collide with enemy
	playerEnemyCollision := Collision.shouldCollide(player, enemy)
	// then
	// ... should allow collision
	if !playerEnemyCollision {
		t.Fatal("Player should collide with enemy")
	}
	// when
	// ... checking if enemy should collide with player
	enemyPlayerCollision := Collision.shouldCollide(enemy, player)
	// then
	// ... should allow collision (bidirectional)
	if !enemyPlayerCollision {
		t.Fatal("Enemy should collide with player (bidirectional)")
	}
	// given
	// ... a health pickup collidable
	healthPickup := &MockCollidable{
		CollisionTags: []string{"health_pickup"},
		ActiveState:   true,
	}
	// when
	// ... checking if player should collide with health pickup
	playerHealthCollision := Collision.shouldCollide(player, healthPickup)
	// then
	// ... should allow collision
	if !playerHealthCollision {
		t.Fatal("Player should collide with health pickup")
	}
	// when
	// ... checking if enemy should collide with health pickup
	enemyHealthCollision := Collision.shouldCollide(enemy, healthPickup)
	// then
	// ... should not allow collision
	if enemyHealthCollision {
		t.Fatal("Enemy should NOT collide with health pickup")
	}
	// when
	// ... checking if object collides with itself
	selfCollision := Collision.shouldCollide(player, player)
	// then
	// ... should not allow self-collision
	if selfCollision {
		t.Fatal("Object should not collide with itself")
	}
}

func TestCollisionDetectionInactive(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... an active player at position (0,0,0)
	// ... an inactive enemy overlapping the player
	// ... both registered with the collision system
	InitCollision()
	activePlayer := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	inactiveEnemy := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"enemy"},
		ActiveState:   false, // Inactive
	}
	Collision.RegisterCollidable(activePlayer)
	Collision.RegisterCollidable(inactiveEnemy)
	// when
	// ... the collision system updates
	Collision.Update()
	// then
	// ... should not trigger collision callback for active player
	// ... and should not trigger collision callback for inactive enemy
	if activePlayer.CallbackCalled {
		t.Fatal("Active player should not receive collision callback from inactive enemy")
	}
	if inactiveEnemy.CallbackCalled {
		t.Fatal("Inactive enemy should not receive collision callback")
	}
}

func TestCollisionDetectionOverlapping(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... a player at position (0,0,0)
	// ... an enemy overlapping the player
	// ... both registered with the collision system
	InitCollision()
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	enemy := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"enemy"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(player)
	Collision.RegisterCollidable(enemy)
	// when
	// ... the collision system updates
	Collision.Update()
	// then
	// ... should trigger collision callback for player
	// ... and should trigger collision callback for enemy
	// ... and player should receive enemy as collision other
	// ... and enemy should receive player as collision other
	if !player.CallbackCalled {
		t.Fatal("Player should receive collision callback")
	}
	if !enemy.CallbackCalled {
		t.Fatal("Enemy should receive collision callback")
	}
	if player.CallbackOther != enemy {
		t.Fatal("Player should receive enemy as collision other")
	}
	if enemy.CallbackOther != player {
		t.Fatal("Enemy should receive player as collision other")
	}
}

func TestCollisionDetectionNonOverlapping(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... a player at position (0,0,0)
	// ... an enemy at position (2,0,0) not overlapping
	// ... both registered with the collision system
	InitCollision()
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	enemy := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 2, Y: 0, Z: 0},
			Max: rl.Vector3{X: 3, Y: 1, Z: 1},
		},
		CollisionTags: []string{"enemy"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(player)
	Collision.RegisterCollidable(enemy)
	// when
	// ... the collision system updates
	Collision.Update()
	// then
	// ... should not trigger collision callback for player
	// ... and should not trigger collision callback for enemy
	if player.CallbackCalled {
		t.Fatal("Player should not receive collision callback for non-overlapping objects")
	}
	if enemy.CallbackCalled {
		t.Fatal("Enemy should not receive collision callback for non-overlapping objects")
	}
}

func TestCheckMovementBlocked(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... a player at position (0,0,0)
	InitCollision()
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	// ... an obstacle touching the player at position (1,0,0)
	// ... both registered with the collision system
	// ... a new position that would cause collision
	obstacle := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 1, Y: 0, Z: 0},
			Max: rl.Vector3{X: 2, Y: 1, Z: 1},
		},
		CollisionTags: []string{"obstacle"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(player)
	Collision.RegisterCollidable(obstacle)
	newPos := rl.Vector3{X: 1.5, Y: 0.5, Z: 0.5}
	// when
	// ... checking if the player can move to the new position
	canMove := Collision.CheckMovement(player, newPos)
	// then
	// ... should block the movement
	if canMove {
		t.Fatal("Movement should be blocked by obstacle")
	}
}

func TestCheckMovementAllowed(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... a player at position (0,0,0)
	// ... an obstacle far away at position (5,0,0)
	// ... both registered with the collision system
	// ... a new position that would not cause collision
	InitCollision()
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	obstacle := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 6, Y: 1, Z: 1},
		},
		CollisionTags: []string{"obstacle"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(player)
	Collision.RegisterCollidable(obstacle)
	newPos := rl.Vector3{X: 2, Y: 0.5, Z: 0.5}
	// when
	// ... checking if the player can move to the new position
	canMove := Collision.CheckMovement(player, newPos)
	// then
	// ... should allow the movement
	if !canMove {
		t.Fatal("Movement should be allowed when no collision would occur")
	}
}

func TestSphereLikeCollision(t *testing.T) {
	// given
	// ... an initialized collision system
	// ... a small bounding box that should be treated as sphere-like
	// ... a large bounding box that should not be treated as sphere-like
	InitCollision()
	sphereBox := rl.BoundingBox{
		Min: rl.Vector3{X: 0, Y: 0, Z: 0},
		Max: rl.Vector3{X: 0.5, Y: 0.5, Z: 0.5},
	}
	largeBox := rl.BoundingBox{
		Min: rl.Vector3{X: 0, Y: 0, Z: 0},
		Max: rl.Vector3{X: 2, Y: 2, Z: 2},
	}
	// when
	// ... checking if the small box is sphere-like
	isSmallSpherelike := Collision.isSpherelike(sphereBox)
	// then
	// ... should detect it as sphere-like
	if !isSmallSpherelike {
		t.Fatal("Small box should be detected as sphere-like")
	}
	// when
	// ... checking if the large box is sphere-like
	isLargeSpherelike := Collision.isSpherelike(largeBox)
	// then
	// ... should not detect it as sphere-like
	if isLargeSpherelike {
		t.Fatal("Large box should not be detected as sphere-like")
	}
}

func TestCollisionSystemEdgeCases(t *testing.T) {
	// given
	// ... an initialized collision system
	InitCollision()
	// when
	// ... the collision system updates with no collidables
	Collision.Update()
	// then
	// ... should not crash
	// given
	// ... a single collidable object
	// ... registered with the collision system
	singleObj := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"test"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(singleObj)
	// when
	// ... the collision system updates with single object
	Collision.Update()
	// then
	// ... should not trigger self-collision
	if singleObj.CallbackCalled {
		t.Fatal("Single object should not trigger collision with itself")
	}
	// given
	// ... a multi-tagged object
	// ... an enemy overlapping the multi-tagged object
	// ... a fresh collision system with both objects
	multiTagObj := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player", "test"},
		ActiveState:   true,
	}
	enemy := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"enemy"},
		ActiveState:   true,
	}
	Collision.ClearAll()
	Collision.RegisterCollidable(multiTagObj)
	Collision.RegisterCollidable(enemy)
	// when
	// ... the collision system updates
	Collision.Update()
	// then
	// ... should detect collision through the player tag
	if !multiTagObj.CallbackCalled {
		t.Fatal("Multi-tag object should collide with enemy")
	}
}

func TestPlayerMovementBlockedByEnemyCollision(t *testing.T) {
	InitCollision()
	// given
	// ... a player at position (0, 0, 0)
	// ... an enemy touching the player at position (1, 0, 0)
	// ... both registered with the collision system
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}

	enemy := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 1, Y: 0, Z: 0},
			Max: rl.Vector3{X: 2, Y: 1, Z: 1},
		},
		CollisionTags: []string{"enemy"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(player)
	Collision.RegisterCollidable(enemy)
	// when
	// ... the collision system updates to detect current collisions
	Collision.Update()
	// then
	// ... both objects should detect the collision
	if !player.CallbackCalled {
		t.Fatal("Player should detect collision with enemy")
	}
	if !enemy.CallbackCalled {
		t.Fatal("Enemy should detect collision with player")
	}
	// when
	// ... the player tries to move away from the enemy (to the left)
	moveLeftPos := rl.Vector3{X: -0.5, Y: 0.5, Z: 0.5}
	canMoveLeft := Collision.CheckMovement(player, moveLeftPos)
	// then
	// ... the player should be able to move away
	if !canMoveLeft {
		t.Fatal("Player should be able to move away from enemy")
	}
	// when
	// ... the player tries to move further into the enemy (to the right)
	moveRightPos := rl.Vector3{X: 1.5, Y: 0.5, Z: 0.5}
	canMoveRight := Collision.CheckMovement(player, moveRightPos)
	// then
	// ... the movement should be blocked
	if canMoveRight {
		t.Fatal("Player should not be able to move into enemy")
	}
	// when
	// ... the player tries to move sideways (perpendicular to collision)
	moveSidewaysPos := rl.Vector3{X: 0.5, Y: 0.5, Z: 1.5}
	canMoveSideways := Collision.CheckMovement(player, moveSidewaysPos)
	// then
	// ... the player should be able to move sideways
	if !canMoveSideways {
		t.Fatal("Player should be able to move perpendicular to enemy collision")
	}
}

func TestPlayerSlideAlongEnemyWhenColliding(t *testing.T) {
	InitCollision()
	// given
	// ... a player at position (0, 0, 0)
	// ... an enemy directly in front of the player
	// ... both registered with the collision system
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}

	enemy := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 1.5, Y: 0, Z: 0},
			Max: rl.Vector3{X: 2.5, Y: 1, Z: 1},
		},
		CollisionTags: []string{"enemy"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(player)
	Collision.RegisterCollidable(enemy)
	// when
	// ... the player tries to move diagonally towards the enemy
	diagonalMove := rl.Vector3{X: 1.5, Y: 0.5, Z: 1.5}
	canMoveDiagonal := Collision.CheckMovement(player, diagonalMove)
	// then
	// ... the diagonal movement should be blocked (moving into enemy)
	if canMoveDiagonal {
		t.Fatal("Player should not be able to move diagonally into enemy")
	}
	// when
	// ... the player tries to move purely sideways (slide along enemy)
	slideMove := rl.Vector3{X: 0.5, Y: 0.5, Z: 1.5}
	canSlide := Collision.CheckMovement(player, slideMove)
	// then
	// ... the slide movement should be allowed
	if !canSlide {
		t.Fatal("Player should be able to slide along enemy")
	}
}

func TestObstacleBlocksAllMovementWhenTouching(t *testing.T) {
	InitCollision()
	// given
	// ... a player at position (0, 0, 0)
	// ... an obstacle touching the player at position (1, 0, 0)
	// ... both registered with the collision system
	player := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 0, Y: 0, Z: 0},
			Max: rl.Vector3{X: 1, Y: 1, Z: 1},
		},
		CollisionTags: []string{"player"},
		ActiveState:   true,
	}
	obstacle := &MockCollidable{
		BoundingBox: rl.BoundingBox{
			Min: rl.Vector3{X: 1, Y: 0, Z: 0},
			Max: rl.Vector3{X: 2, Y: 1, Z: 1},
		},
		CollisionTags: []string{"obstacle"},
		ActiveState:   true,
	}
	Collision.RegisterCollidable(player)
	Collision.RegisterCollidable(obstacle)
	// when
	// ... the player tries to move sideways (perpendicular to collision)
	moveSidewaysPos := rl.Vector3{X: 0.5, Y: 0.5, Z: 1.5}
	canMoveSideways := Collision.CheckMovement(player, moveSidewaysPos)
	// then
	// ... the movement should be blocked (obstacles block all movement when touching)
	if canMoveSideways {
		t.Fatal("Player should NOT be able to slide along obstacles")
	}
	// when
	// ... the player tries to move away from the obstacle
	moveAwayPos := rl.Vector3{X: -0.5, Y: 0.5, Z: 0.5}
	canMoveAway := Collision.CheckMovement(player, moveAwayPos)
	// then
	// ... the movement away should be allowed
	if !canMoveAway {
		t.Fatal("Player should be able to move away from obstacle")
	}
}
