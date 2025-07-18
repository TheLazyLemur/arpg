package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Player represents the player character
type Player struct {
	position rl.Vector3
	rotation float32
	speed    float32
}

// Camera represents the game camera
type Camera struct {
	camera rl.Camera3D
	offset rl.Vector3
}

// Obstacle represents a collidable object
type Obstacle struct {
	position rl.Vector3
	size     rl.Vector3
	radius   float32
	height   float32
	isBox    bool
}

// Bullet represents a projectile
type Bullet struct {
	position rl.Vector3
	velocity rl.Vector3
	lifetime float32
	speed    float32
	radius   float32
	active   bool
}

// Enemy represents a target to shoot
type Enemy struct {
	position  rl.Vector3
	health    float32
	maxHealth float32
	radius    float32
	height    float32
	active    bool
}

// Game state to hold obstacles, bullets, and enemies
type GameState struct {
	obstacles []Obstacle
	bullets   []Bullet
	enemies   []Enemy
}
