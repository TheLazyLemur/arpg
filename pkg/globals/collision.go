package globals

import (
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CollisionSystem struct {
	collidables []Collidable
}

var Collision *CollisionSystem

func InitCollision() {
	Collision = &CollisionSystem{
		collidables: make([]Collidable, 0),
	}
}

type Collidable interface {
	GetBoundingBox() rl.BoundingBox
	GetCollisionTags() []string
	OnCollision(other Collidable)
	IsActive() bool
}

func (cs *CollisionSystem) RegisterCollidable(obj Collidable) {
	cs.collidables = append(cs.collidables, obj)
}

func (cs *CollisionSystem) UnregisterCollidable(obj Collidable) {
	for i, collidable := range cs.collidables {
		if collidable == obj {
			cs.collidables = slices.Delete(cs.collidables, i, i+1)
			return
		}
	}
}

// ClearAll removes all collidables from the collision system
func (cs *CollisionSystem) ClearAll() {
	cs.collidables = make([]Collidable, 0)
}

func (cs *CollisionSystem) Update() {
	for i := range cs.collidables {
		objA := cs.collidables[i]
		if !objA.IsActive() {
			continue
		}

		for j := i + 1; j < len(cs.collidables); j++ {
			objB := cs.collidables[j]
			if !objB.IsActive() {
				continue
			}

			if cs.shouldCollide(objA, objB) {
				if cs.checkCollision(objA, objB) {
					// Trigger collision callbacks for both objects
					objA.OnCollision(objB)
					objB.OnCollision(objA)
				}
			}
		}
	}
}

func (cs *CollisionSystem) shouldCollide(objA, objB Collidable) bool {
	if objA == objB {
		return false
	}

	if !objA.IsActive() || !objB.IsActive() {
		return false
	}

	tagsA := objA.GetCollisionTags()
	tagsB := objB.GetCollisionTags()

	collisionRules := map[string][]string{
		"player":        {"obstacle", "enemy", "health_pickup"},
		"bullet":        {"obstacle", "enemy"},
		"enemy":         {"player", "bullet", "obstacle"},
		"obstacle":      {"player", "bullet", "enemy", "obstacle"},
		"health_pickup": {"player"},
	}

	for _, tagA := range tagsA {
		if allowedTags, exists := collisionRules[tagA]; exists {
			for _, tagB := range tagsB {
				if slices.Contains(allowedTags, tagB) {
					return true
				}
			}
		}
	}

	return false
}

func (cs *CollisionSystem) checkCollision(objA, objB Collidable) bool {
	boxA := objA.GetBoundingBox()
	boxB := objB.GetBoundingBox()

	if cs.isSpherelike(boxA) || cs.isSpherelike(boxB) {
		if cs.isSpherelike(boxA) {
			center := rl.Vector3{
				X: (boxA.Min.X + boxA.Max.X) / 2,
				Y: (boxA.Min.Y + boxA.Max.Y) / 2,
				Z: (boxA.Min.Z + boxA.Max.Z) / 2,
			}
			radius := (boxA.Max.X - boxA.Min.X) / 2
			return rl.CheckCollisionBoxSphere(boxB, center, radius)
		} else {
			center := rl.Vector3{
				X: (boxB.Min.X + boxB.Max.X) / 2,
				Y: (boxB.Min.Y + boxB.Max.Y) / 2,
				Z: (boxB.Min.Z + boxB.Max.Z) / 2,
			}
			radius := (boxB.Max.X - boxB.Min.X) / 2
			return rl.CheckCollisionBoxSphere(boxA, center, radius)
		}
	}

	return rl.CheckCollisionBoxes(boxA, boxB)
}

func (cs *CollisionSystem) isSpherelike(box rl.BoundingBox) bool {
	width := box.Max.X - box.Min.X
	height := box.Max.Y - box.Min.Y
	depth := box.Max.Z - box.Min.Z

	return width < 1.0 && height < 1.0 && depth < 1.0
}

func (cs *CollisionSystem) CheckMovement(obj Collidable, newPos rl.Vector3) bool {
	originalBox := obj.GetBoundingBox()
	offset := rl.Vector3{
		X: newPos.X - (originalBox.Min.X+originalBox.Max.X)/2,
		Y: newPos.Y - (originalBox.Min.Y+originalBox.Max.Y)/2,
		Z: newPos.Z - (originalBox.Min.Z+originalBox.Max.Z)/2,
	}

	tempBox := rl.BoundingBox{
		Min: rl.Vector3{
			X: originalBox.Min.X + offset.X,
			Y: originalBox.Min.Y + offset.Y,
			Z: originalBox.Min.Z + offset.Z,
		},
		Max: rl.Vector3{
			X: originalBox.Max.X + offset.X,
			Y: originalBox.Max.Y + offset.Y,
			Z: originalBox.Max.Z + offset.Z,
		},
	}

	for _, other := range cs.collidables {
		if other == obj || !other.IsActive() {
			continue
		}

		if cs.shouldCollide(obj, other) {
			otherBox := other.GetBoundingBox()
			if rl.CheckCollisionBoxes(tempBox, otherBox) {
				otherTags := other.GetCollisionTags()
				isEnemy := slices.Contains(otherTags, "enemy")

				if isEnemy && rl.CheckCollisionBoxes(originalBox, otherBox) {
					if !cs.isMovementIncreasingPenetration(originalBox, tempBox, otherBox) {
						continue // Allow this movement
					}
				}

				return false
			}
		}
	}

	return true
}

func (cs *CollisionSystem) isMovementIncreasingPenetration(
	originalBox, newBox, targetBox rl.BoundingBox,
) bool {
	originalPenetrationX := cs.getPenetrationDepth(
		originalBox.Min.X, originalBox.Max.X,
		targetBox.Min.X, targetBox.Max.X,
	)
	originalPenetrationZ := cs.getPenetrationDepth(
		originalBox.Min.Z, originalBox.Max.Z,
		targetBox.Min.Z, targetBox.Max.Z,
	)

	newPenetrationX := cs.getPenetrationDepth(
		newBox.Min.X, newBox.Max.X,
		targetBox.Min.X, targetBox.Max.X,
	)
	newPenetrationZ := cs.getPenetrationDepth(
		newBox.Min.Z, newBox.Max.Z,
		targetBox.Min.Z, targetBox.Max.Z,
	)

	return newPenetrationX > originalPenetrationX || newPenetrationZ > originalPenetrationZ
}

func (cs *CollisionSystem) getPenetrationDepth(min1, max1, min2, max2 float32) float32 {
	if min1 > max2 || max1 < min2 {
		return 0 // No overlap
	}

	overlapStart := max(min1, min2)
	overlapEnd := min(max1, max2)
	return overlapEnd - overlapStart
}
