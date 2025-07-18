package globals

import (
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TriggerSystem struct {
	triggers    []Triggerable
	collidables []Collidable
}

var Triggers *TriggerSystem

func InitTriggers() {
	Triggers = &TriggerSystem{
		triggers:    make([]Triggerable, 0),
		collidables: make([]Collidable, 0),
	}
}

type Triggerable interface {
	GetTriggerBounds() rl.BoundingBox
	GetTriggerTags() []string
	OnTriggerEnter(other Collidable)
	IsActive() bool
}

func (ts *TriggerSystem) RegisterTrigger(trigger Triggerable) {
	ts.triggers = append(ts.triggers, trigger)
}

func (ts *TriggerSystem) UnregisterTrigger(trigger Triggerable) {
	for i, t := range ts.triggers {
		if t == trigger {
			ts.triggers = slices.Delete(ts.triggers, i, i+1)
			return
		}
	}
}

func (ts *TriggerSystem) RegisterCollidable(collidable Collidable) {
	ts.collidables = append(ts.collidables, collidable)
}

func (ts *TriggerSystem) UnregisterCollidable(collidable Collidable) {
	for i, c := range ts.collidables {
		if c == collidable {
			ts.collidables = slices.Delete(ts.collidables, i, i+1)
			return
		}
	}
}

func (ts *TriggerSystem) ClearAll() {
	ts.triggers = make([]Triggerable, 0)
	ts.collidables = make([]Collidable, 0)
}

func (ts *TriggerSystem) Update() {
	for _, trigger := range ts.triggers {
		if !trigger.IsActive() {
			continue
		}

		triggerBounds := trigger.GetTriggerBounds()

		for _, collidable := range ts.collidables {
			if !collidable.IsActive() {
				continue
			}

			if ts.shouldTrigger(trigger, collidable) {
				if ts.checkTriggerCollision(triggerBounds, collidable.GetBoundingBox()) {
					trigger.OnTriggerEnter(collidable)
				}
			}
		}
	}
}

func (ts *TriggerSystem) shouldTrigger(trigger Triggerable, collidable Collidable) bool {
	triggerTags := trigger.GetTriggerTags()
	collidableTags := collidable.GetCollisionTags()

	triggerRules := map[string][]string{
		"health_pickup": {"player"},
	}

	for _, triggerTag := range triggerTags {
		if allowedTags, exists := triggerRules[triggerTag]; exists {
			for _, collidableTag := range collidableTags {
				if slices.Contains(allowedTags, collidableTag) {
					return true
				}
			}
		}
	}

	return false
}

func (ts *TriggerSystem) checkTriggerCollision(
	triggerBounds, collidableBounds rl.BoundingBox,
) bool {
	return rl.CheckCollisionBoxes(triggerBounds, collidableBounds)
}
