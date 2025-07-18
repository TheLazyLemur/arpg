package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// ActionType represents different types of actions entities can request
type ActionType int

const (
	ActionTypeSpawnBullet ActionType = iota
)

// Action represents an action that an entity wants to perform
type Action struct {
	Type ActionType
	Data interface{}
}

// SpawnBulletData contains data for spawning a bullet
type SpawnBulletData struct {
	Position  rl.Vector3
	Direction rl.Vector3
	Speed     float32
	Lifetime  float32
	Damage    float32
}

// ActionHandler interface for handling entity actions
type ActionHandler interface {
	HandleAction(action Action)
}

